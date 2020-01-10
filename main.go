package main

import (
	"context"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
)

var clients = &sync.Map{}
var broadcast = make(chan Data, 100000)
var upgrader = websocket.Upgrader{}
var onlineAtom int64
var pool = &sync.Pool{New: func() interface{} {
	return make([]byte, 0, 512)
}}

type Data struct {
	mType   int
	message []byte
	from    *websocket.Conn
}

func IncAtom(i *int64) {
	atomic.AddInt64(i, 1)
}

func DecrAtom(i *int64) {
	atomic.AddInt64(i, -1)
}

func LoadAtom(i *int64) int64 {
	return atomic.LoadInt64(i)
}

//handleMessages
func handleMessages() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			handleMessages()
		}
	}()

	for msg := range broadcast {
		clients.Range(func(key, value interface{}) bool {
			key.(chan Data) <- msg
			return true
		})
	}
}

//writer
func writer(ctx context.Context, ws *websocket.Conn) {
	var message = make(chan Data, 100000)
	defer func() {
		clients.Delete(message)
		close(message)
	}()
	ctx, cancel := context.WithCancel(ctx)
	clients.Store(message, cancel)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-message:
			if msg.from != ws {
				writer, err := ws.NextWriter(msg.mType)
				if err != nil {
					log.Println(ws.RemoteAddr().String(), err)
					return
				}
				if _, err := writer.Write(msg.message); err != nil {
					log.Println(ws.RemoteAddr().String(), err)
					return
				}
			}
		}
	}
}

//handler
func handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println(ws.RemoteAddr().String(), "accept")
	IncAtom(&onlineAtom)
	defer func() {
		DecrAtom(&onlineAtom)
		ws.Close()
	}()

	go writer(context.Background(), ws)

	for {
		mt, reader, err := ws.NextReader()
		if err != nil {
			log.Println(ws.RemoteAddr().String(), err)
			return
		}

		buf := pool.Get().([]byte)
		//насыщаем буфер по размеру капасити, небольшой хак т.к io.Read не умеет в апенды
		buf = buf[:cap(buf)]
		n, err := reader.Read(buf)
		if err != nil && err != io.ErrUnexpectedEOF {
			log.Println(ws.RemoteAddr().String(), err)
			return
		}

		broadcast <- Data{from: ws, mType: mt, message: buf[:n]}

		//очищаем буфер
		buf = buf[:0]
		pool.Put(buf)

		runtime.Gosched()
	}
}

func main() {
	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handler)
	go handleMessages()
	log.Println("run")
	log.Fatal(http.ListenAndServe(":80", nil))
}

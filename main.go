package main

import (
	"context"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var clients = &sync.Map{}
var broadcast = make(chan Data, 100000)
var upgrader = websocket.Upgrader{}

type Data struct {
	mType   int
	message []byte
	from    *websocket.Conn
}

func handleMessages() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	for {
		msg := <-broadcast
		clients.Range(func(key, value interface{}) bool {
			key.(chan Data) <- msg
			return true
		})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	log.Println(ws.RemoteAddr().String(), "accept")
	var message = make(chan Data, 100000)
	defer func() {
		ws.Close()
		clients.Delete(message)
		close(message)
	}()

	ctx, cancel := context.WithCancel(context.Background())
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
		default:
			var d Data
			mtype, reader, err := ws.NextReader()
			if err != nil {
				log.Println(ws.RemoteAddr().String(), err)
				return
			}

			buf, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Println(ws.RemoteAddr().String(), err)
				continue
			}
			d.from = ws
			d.mType = mtype
			d.message = buf
			broadcast <- d
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handler)
	go handleMessages()
	log.Println("run")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

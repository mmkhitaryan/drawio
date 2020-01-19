package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const quota = 1800
const version = "0.0.1"

var clients = &sync.Map{}
var broadcast = make(chan *Data, 100000)

//TODO временый костыйль с origin
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}
var onlineAtom int64
var pool = &sync.Pool{New: func() interface{} {
	return make([]byte, 0, 512*5)
}}

//Data
type Data struct {
	mType   int
	force   bool
	message []byte
	from    *websocket.Conn
}

//socket
type socket struct {
	conn       *websocket.Conn
	message    chan Data
	quotum     int64
	quotumLock uint32
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
			key.(chan Data) <- *msg
			return true
		})
	}
}

//writer
func (s *socket) writer(ctx context.Context) {
	timer := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if s.statusQuotum() == 1 {
				i := atomic.AddInt64(&s.quotum, -600)
				s.sendQuotum(i, 1)
				if i <= 0 {
					s.unlockQuotum()
				}
			}
		case msg := <-s.message:
			if msg.from != s.conn || msg.force {
				writer, err := s.conn.NextWriter(msg.mType)
				if err != nil {
					poolPut(msg.message)
					log.Println(s.conn.RemoteAddr().String(), err)
					return
				}
				if _, err := writer.Write(msg.message); err != nil {
					poolPut(msg.message)
					log.Println(s.conn.RemoteAddr().String(), err)
					return
				}

				if err = writer.Close(); err != nil {
					poolPut(msg.message)
					log.Println(s.conn.RemoteAddr().String(), err)
					return
				}
				poolPut(msg.message)
			}
		}
	}
}

func poolPut(buf []byte) {
	buf = buf[:0]
	pool.Put(buf)
}

//socketPayload
type socketPayload struct {
	ID   int `json:"id"`
	Data struct {
		Points [][]int `json:"points,omitempty"`
	} `json:"data,omitempty"`
	Count int64 `json:"count,omitempty"`
}

//reader
func (s *socket) reader() {
	ctx, cancel := context.WithCancel(context.Background())
	atomic.AddInt64(&onlineAtom, 1)
	defer func() {
		cancel()
		atomic.AddInt64(&onlineAtom, -1)
		clients.Delete(s.message)
		close(s.message)
		s.conn.Close()
	}()
	clients.Store(s.message, struct{}{})
	go s.writer(ctx)

	for {
		mt, reader, err := s.conn.NextReader()
		if err != nil {
			log.Println(s.conn.RemoteAddr().String(), err)
			return
		}

		buf := pool.Get().([]byte)
		//насыщаем буфер по размеру капасити, небольшой хак т.к io.Read не умеет в апенды
		buf = buf[:cap(buf)]
		n, err := io.ReadFull(reader, buf)
		if err != nil && err != io.ErrUnexpectedEOF {
			log.Println(s.conn.RemoteAddr().String(), err)
			return
		}
		var sp socketPayload
		if err := json.Unmarshal(buf[:n], &sp); err != nil {
			log.Println(string(buf[:n]))
			log.Println(err)
		}

		if countQuota, ok := s.quotumAllow(len(sp.Data.Points)); !ok {
			s.sendQuotum(countQuota, mt)
			buf = buf[:0]
			pool.Put(buf)
			continue
		}
		broadcast <- &Data{from: s.conn, mType: mt, force: false, message: buf[:n]}
	}
}

//quotumAllow
func (s *socket) quotumAllow(count int) (int64, bool) {
	if s.statusQuotum() == 1 {
		return atomic.LoadInt64(&s.quotum), false
	}
	i := atomic.LoadInt64(&s.quotum)
	if i >= quota {
		s.lockQuotum()
		return i, false
	}
	i = atomic.AddInt64(&s.quotum, int64(count))
	return i, true
}

func (s *socket) lockQuotum() {
	if atomic.LoadUint32(&s.quotumLock) == 1 {
		return
	}
	atomic.StoreUint32(&s.quotumLock, 1)
}

func (s *socket) unlockQuotum() {
	if atomic.LoadUint32(&s.quotumLock) == 1 {
		atomic.StoreUint32(&s.quotumLock, 0)
	}
}

func (s *socket) statusQuotum() uint32 {
	return atomic.LoadUint32(&s.quotumLock)
}

//sendQuotum
func (s *socket) sendQuotum(quotum int64, mt int) {
	sp := &socketPayload{
		ID:    2,
		Count: quotum,
	}
	data, _ := json.Marshal(sp)
	s.message <- Data{from: s.conn, mType: mt, force: true, message: data}
}

//upgrade
func upgrade(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	s := &socket{conn: ws, message: make(chan Data, 100000)}
	go s.reader()
}

//online
func online(w http.ResponseWriter, r *http.Request) {
	total := atomic.LoadInt64(&onlineAtom)
	fmt.Fprint(w, total)
}

//printVersion
func printVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, version)
}

func main() {
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	http.HandleFunc("/online", online)
	http.HandleFunc("/version", printVersion)
	http.HandleFunc("/ws", upgrade)
	go handleMessages()
	log.Println("run")
	log.Fatal(http.ListenAndServe(":80", nil))
}

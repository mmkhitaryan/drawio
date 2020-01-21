package main

import (
	"container/ring"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const quota = 1800

var clients = &sync.Map{}
var online int64
var broadcast = make(chan *data, 100000)
var circularBuf = &safeCircularBuffer{mu: &sync.RWMutex{}, rbuf: ring.New(1000)}

type safeCircularBuffer struct {
	mu   *sync.RWMutex
	rbuf *ring.Ring
}

//data
type data struct {
	mType   int
	force   bool
	message []byte
	from    *websocket.Conn
}

//socket
type socket struct {
	conn       *websocket.Conn
	message    chan data
	quotum     int64
	quotumLock uint32
}

//TODO временый костыйль с origin
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

//handleMessages
func handleMessages() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		go handleMessages()
	}()

	for msg := range broadcast {
		clients.Range(func(key, value interface{}) bool {
			key.(chan data) <- *msg
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
			if s.quotumIsLock() {
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
					s.log(err)
					return
				}
				if _, err := writer.Write(msg.message); err != nil {
					s.log(err)
					return
				}

				if err = writer.Close(); err != nil {
					s.log(err)
					return
				}
			}
		}
	}
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
	atomic.AddInt64(&online, 1)
	defer func() {
		cancel()
		atomic.AddInt64(&online, -1)
		clients.Delete(s.message)
		close(s.message)
		s.conn.Close()
	}()
	s.sendHistory()
	clients.Store(s.message, struct{}{})
	go s.writer(ctx)

	for {
		mt, reader, err := s.conn.NextReader()
		if err != nil {
			s.log(err)
			return
		}

		buf, err := ioutil.ReadAll(reader)
		if err != nil && err != io.ErrUnexpectedEOF {
			s.log(err)
			return
		}

		var sp socketPayload
		if err := json.Unmarshal(buf, &sp); err != nil {
			s.log(err)
		}

		if countQuota, ok := s.quotumAllow(len(sp.Data.Points)); !ok {
			s.sendQuotum(countQuota, mt)
			continue
		}
		broadcast <- &data{from: s.conn, mType: mt, force: false, message: buf}
		circularBuf.mu.Lock()
		circularBuf.rbuf.Value = buf
		circularBuf.rbuf = circularBuf.rbuf.Next()
		circularBuf.mu.Unlock()
	}
}

func (s *socket) sendHistory() {
	circularBuf.mu.RLock()
	defer circularBuf.mu.RUnlock()
	circularBuf.rbuf.Do(func(val interface{}) {
		if msg, ok := val.([]byte); ok {
			s.message <- data{from: s.conn, mType: 1, force: true, message: msg}
		}
	})
}

//quotumAllow
func (s *socket) quotumAllow(count int) (int64, bool) {
	if s.quotumIsLock() {
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
	if s.quotumIsLock() {
		return
	}
	atomic.StoreUint32(&s.quotumLock, 1)
}

func (s *socket) unlockQuotum() {
	if s.quotumIsLock() {
		atomic.StoreUint32(&s.quotumLock, 0)
	}
}

func (s *socket) quotumIsLock() bool {
	return atomic.LoadUint32(&s.quotumLock) == 1
}

//sendQuotum
func (s *socket) sendQuotum(quotum int64, mt int) {
	sp := &socketPayload{
		ID:    2,
		Count: quotum,
	}
	if msg, err := json.Marshal(sp); err == nil {
		s.message <- data{from: s.conn, mType: mt, force: true, message: msg}
	}
}

func (s *socket) log(err error) {
	log.Println(s.conn.RemoteAddr().String(), err)
}

//upgrade
func upgrade(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	s := &socket{conn: ws, message: make(chan data, 100000)}
	go s.reader()
}

//online
func onlineHandle(w http.ResponseWriter, r *http.Request) {
	total := atomic.LoadInt64(&online)
	fmt.Fprint(w, total)
}

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	go handleMessages()

	http.HandleFunc("/online", onlineHandle)
	http.HandleFunc("/ws", upgrade)
	log.Fatal(http.ListenAndServe(":80", nil))
}

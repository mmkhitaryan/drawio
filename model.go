package main

import (
	"golang.org/x/net/websocket"
	"sync"
)

type ClientsPool struct {
	sync.Mutex
	clients map[*websocket.Conn]struct{}
	Count   int
}

func NewClientsPool() ClientsPool {
	return ClientsPool{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (c *ClientsPool) Add(ws *websocket.Conn) {
	c.Lock()
	defer c.Unlock()

	c.Count++
	c.clients[ws] = struct{}{}
}

func (c *ClientsPool) GetAll() []*websocket.Conn {
	c.Lock()
	defer c.Unlock()

	all := make([]*websocket.Conn, 0, len(c.clients))
	for k := range c.clients {
		all = append(all, k)
	}
	return all
}

func (c *ClientsPool) Del(ws *websocket.Conn) {
	c.Lock()
	defer c.Unlock()

	err := ws.Close()
	if err != nil {
		println("[E] ws.Close()", err)
	}
	c.Count--
	delete(c.clients, ws)
}

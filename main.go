package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	//"github.com/gorilla/websocket"
)

var clients = NewClientsPool()           // connected clients
var broadcast = make(chan Message, 1000) // broadcast channel

type Message struct {
	OldX int `json:"old_x"`
	OldY int `json:"old_y"`
	NewX int `json:"new_x"`
	NewY int `json:"new_y"`
	conn *websocket.Conn
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		body, _ := json.Marshal(msg)
		for _, client := range clients.GetAll() {
			go func(client *websocket.Conn) {
				if client != msg.conn {
					n, err := client.Write(body)
					if err != nil || n == 0 {
						log.Printf("error: %v", err)
						clients.Del(msg.conn)
					}
				}
			}(client)
		}
	}
}

func handler(ws *websocket.Conn) {

	clients.Add(ws)

	for {
		var msg Message
		msg.conn = ws // Also include sener's conn object
		// Read in a new message as JSON and map it to a Message object

		dec := json.NewDecoder(ws)
		err := dec.Decode(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			clients.Del(ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func main() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.Handle("/ws", websocket.Handler(handler))
	go handleMessages()
	log.Fatal(http.ListenAndServe(":80", nil))
}

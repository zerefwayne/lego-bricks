package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

func handleMessages() {

	for {

		msg := <-broadcast

		for client := range clients {

			err := client.WriteJSON(msg)

			if err != nil {
				log.Printf("error: %+v\n", err)
				client.Close()
				delete(clients, client)
			}

		}

	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	print("New connection")

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {

		var msg Message

		print("New message")

		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %+v\n", err)
			delete(clients, ws)
			break
		}

		broadcast <- msg

	}

}

func main() {

	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("http server starting on port :8000")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}

}

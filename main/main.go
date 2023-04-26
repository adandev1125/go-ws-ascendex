package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	// create a Dialer
	dialer := websocket.DefaultDialer
	// set up the WebSocket URL
	url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	// dial the WebSocket server
	conn, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server: ", err)
	}
}

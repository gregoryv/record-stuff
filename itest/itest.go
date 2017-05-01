package main

import (
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:8080/recordings/1234"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Streaming...")
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	log.Print("Done")
}

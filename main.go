package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"io/ioutil"
)

func recordHandler(ws *websocket.Conn) {
	var data []byte
	websocket.Message.Receive(ws, &data)
	ioutil.WriteFile("/tmp/out.wav", data, 0777)
	ws.Write([]byte("saved"))
	// to keep websocket open you cannot return here
}

func main() {
	http.Handle("/record", websocket.Handler(recordHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	// getUserMedia will not work on insecure origins
	// https://goo.gl/Y0ZkNV
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

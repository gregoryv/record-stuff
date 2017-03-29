package main

import (
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"log"
)

func recordHandler(ws *websocket.Conn) {
	var data []byte
	websocket.Message.Receive(ws, &data)
	ioutil.WriteFile("/tmp/out.wav", data, 0777)
	ws.Write([]byte("saved"))
	// to keep websocket open you cannot return here
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("soundBlob")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	out := "/tmp/" + r.FormValue("filename")
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	log.Printf("%v", out)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/record", websocket.Handler(recordHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	// getUserMedia will not work on insecure origins
	// https://goo.gl/Y0ZkNV
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

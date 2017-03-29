package main

import (
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"log"
	"github.com/gorilla/mux"
	"time"
	"fmt"
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

func homeHandler(w http.ResponseWriter, r *http.Request)  {
	fh, err := os.Open("static/index.html")
	if err != nil {
		log.Print(err)
		fmt.Fprint(w, "404 not found")
		return
	}
	defer fh.Close()
	io.Copy(w, fh)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.HandleFunc("/upload", uploadHandler)
	r.Handle("/record", websocket.Handler(recordHandler))
    http.Handle("/", r)
	
	// getUserMedia will not work on insecure origins
	// https://goo.gl/Y0ZkNV
	
	srv := &http.Server{
        Handler:      r,
        Addr:         ":8080",
		// Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

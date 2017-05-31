// Web service for recording streamed audio
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gregoryv/service-api"
	"log"
	"net/http"
	"time"
)

var bind = flag.String("bind", "localhost:8080", "Listen on host:port, must be accessible from inet")

func main() {
	flag.Parse()

	// Register routes
	r := mux.NewRouter()
	initRecordHandlers(r)
	//initUploadHandlers(r)

	// API documentation route
	r.HandleFunc("/", api.WriteServiceSpec)

	// Static content
	upath := `/static/{rest:[a-zA-Z0-9=\-\/\.]*}`
	r.HandleFunc(upath, WriteAsset)
	api.Doc("GET", upath, "Returns demonstration app")

	// Combine all
	http.Handle("/", r)

	// Define server
	srv := &http.Server{
		Handler:      r,
		Addr:         *bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Listening on %s", *bind)
	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func WriteAsset(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.String()
	if upath == "/static/" {
		upath = "/static/index.html"
	}
	log.Print(upath[1:])
	data, err := Asset(upath[1:])
	if err != nil {
		log.Print("not found")
	}
	if upath == "/static/js/main.js" {
		data = bytes.Replace(data, []byte("localhost:8080"), []byte(*bind), 1)
	}
	fmt.Fprint(w, string(data))
}

/*

Notes!

getUserMedia will not work on insecure origins
https://goo.gl/Y0ZkNV

*/

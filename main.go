// Web service for recording streamed audio
package main

import (
	"github.com/gorilla/mux"
	"github.com/gregoryv/service-api"
	"log"
	"net/http"
	"time"
)

func main() {
	// Register routes
	r := mux.NewRouter()
	initRecordHandlers(r)
	//initUploadHandlers(r)

	// API documentation route
	r.HandleFunc("/", api.WriteServiceSpec)

	// Static content
	upath := "/static/"
	r.PathPrefix(upath).Handler(http.FileServer(http.Dir(".")))
	api.Doc("GET", upath, "Returns demonstration app")

	// Combine all
	http.Handle("/", r)

	// Define server
	bind := ":8080"
	srv := &http.Server{
		Handler:      r,
		Addr:         bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Listening on %s", bind)
	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

/*

Notes!

getUserMedia will not work on insecure origins
https://goo.gl/Y0ZkNV

*/

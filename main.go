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
	r := mux.NewRouter()
	initRecordHandlers(r)
	initUploadHandlers(r)
	r.HandleFunc("/", api.WriteServiceSpec)

	upath := "/static/"
	r.PathPrefix(upath).Handler(http.FileServer(http.Dir(".")))
	api.Doc("GET", upath, "Returns demonstration app")

	http.Handle("/", r)

	// getUserMedia will not work on insecure origins
	// https://goo.gl/Y0ZkNV
	bind := ":8080"
	srv := &http.Server{
		Handler: r,
		Addr:    bind,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Listening on %s", bind)
	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gregoryv/service-api"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	r := mux.NewRouter()
	initRecordHandlers(r)
	r.HandleFunc("/", api.WriteServiceSpec)
	r.HandleFunc("/app", WriteDemoApp)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.HandleFunc("/upload", uploadHandler)
	http.Handle("/", r)

	// getUserMedia will not work on insecure origins
	// https://goo.gl/Y0ZkNV
	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("soundBlob")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	out := path.Join(OUT, r.FormValue("filename"))
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	log.Printf("%v", out)
}

func WriteDemoApp(w http.ResponseWriter, r *http.Request) {
	fh, err := os.Open("static/index.html")
	if err != nil {
		log.Print(err)
		fmt.Fprint(w, "404 not found")
		return
	}
	defer fh.Close()
	io.Copy(w, fh)
}

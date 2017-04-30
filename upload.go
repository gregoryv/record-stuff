package main

import (
	"github.com/gorilla/mux"
	"github.com/gregoryv/service-api"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func initUploadHandlers(r *mux.Router) {
	upath := "/upload"
	r.HandleFunc(upath, uploadHandler)
	api.Doc("POST", upath, "Upload recording")
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

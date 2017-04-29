package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	// Where the recordings are stored
	OUT = "/tmp/recordings"
)

// init makes sure there is a place for the recordings
func init() {
	err := os.Mkdir(OUT, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

// The stream from a client will come in chunks
func recordHandler(ws *websocket.Conn) {
	log.Print("Connected")
	// Connection ends if nothing is received if deadline is reached
	deadline := time.Now().Add(time.Second * 5)
	ws.SetDeadline(deadline)

	tofile := path.Join(OUT, "rec.wav")
	file, err := os.Create(tofile)
	if err != nil {
		panic(err)
	}

	for {
		var data []byte
		websocket.Message.Receive(ws, &data)
		if len(data) == 0 {
			// client closed connection
			break
		}
		n, err := file.Write(data)
		if err != nil {
			log.Print(err)
		}
		log.Printf("Saved %d bytes to %s", n, tofile)
	}
	file.Close()
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

type Rec struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

// listRecordings writes json array of recordings
func listRecordings(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(OUT)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Listing")
	names := make([]Rec, len(files))
	for i, file := range files {
		names[i] = Rec{
			Href: "/recordings/" + file.Name(),
			Name: file.Name(),
		}
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(names)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/app", WriteDemoApp)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/recordings/{name}.wav").Handler(
		http.StripPrefix(
			"/recordings/",
			http.FileServer(http.Dir(OUT)),
		),
	).Methods("GET")
	r.HandleFunc("/recordings/", listRecordings).Methods("GET")
	r.HandleFunc("/upload", uploadHandler)
	// Its up to the client to decide where the recording is saved
	// TODO maybe better to use /recordings with different schema or eg. without .wav
	r.Handle("/record", websocket.Handler(recordHandler))
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

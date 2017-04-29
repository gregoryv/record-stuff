package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gregoryv/service-api"
	"golang.org/x/net/websocket"
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

type Rec struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

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

func initRecordHandlers(r *mux.Router) {
	upath := "/recordings/{name}.wav"
	r.PathPrefix(upath).Handler(
		http.StripPrefix(
			"/recordings/",
			http.FileServer(http.Dir(OUT)),
		),
	).Methods("GET")
	api.Doc("GET", upath, "Returns the recorded audio")

	upath = "/recordings/"
	r.HandleFunc(upath, listRecordings).Methods("GET")
	api.Doc("GET", upath, "Returns list of all recordings").Resource = "Rec"

	// Its up to the client to decide where the recording is saved
	// TODO maybe better to use /recordings with different schema or eg. without .wav
	r.Handle("/record", websocket.Handler(recordHandler))

}

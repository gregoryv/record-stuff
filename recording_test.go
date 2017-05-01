package main

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func Test_recordHandler(t *testing.T) {
	// Setup service under test
	r := mux.NewRouter()
	initRecordHandlers(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	filename := "hi"
	expectedPath := path.Join(OUT, filename)
	// Clean up first
	if err := os.RemoveAll(expectedPath); err != nil {
		t.Fatalf("%s", err)
	}

	// Fix test url
	turl := strings.Replace(ts.URL, "http", "ws", 1) + "/recordings/" + filename
	origin := "http://localhost/"
	ws, err := websocket.Dial(turl, "", origin)
	if err != nil {
		t.Fatalf("%s", err)
	}
	// Simulate streaming
	if _, err := ws.Write([]byte("audio data here")); err != nil {
		t.Errorf("%s", err)
	}
	ts.Close()                   // Finnish stream by closing connection
	time.Sleep(time.Millisecond) // Give the service a moment to save the file
	if _, err = os.Open(expectedPath); err != nil {
		t.Errorf("%s", err)
	}
}

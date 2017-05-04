package main

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"net/http/httptest"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_recordHandler(t *testing.T) {
	// Setup service under test
	r := mux.NewRouter()
	initSocketHandler(r) // private func, should not be used in tests
	ts := httptest.NewServer(r)
	defer ts.Close()
	// We know where a file should be saved
	before, err := filepath.Glob(path.Join(OUT, "*"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	// Fix test url
	turl := strings.Replace(ts.URL, "http", "ws", 1) + "/recordings/test"
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
	// Count files
	after, err := filepath.Glob(path.Join(OUT, "*"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	// Assert
	if len(after)-1 != len(before) {
		t.Error("There should be one more file after upload")
	}
}

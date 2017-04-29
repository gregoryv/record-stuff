package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(WriteDemoApp))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Contains(body, []byte("<html>")) {
		t.Errorf("HomeHandler should return html")
	}
}

func ExampleWriteServiceSpec() {
	ts := httptest.NewServer(http.HandlerFunc(WriteServiceSpec))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
	// Output:
	// {}
}

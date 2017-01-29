package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

var (
	server *httptest.Server
	endpoint string
)

func init() {
	server = httptest.NewServer(AddHandlers())
	endpoint = fmt.Sprintf("%s/status", server.URL)
	d := LoadDataFile("data.json")
	fmt.Println(d)
}

func TestviewCountryCodes(t *testing.T) {
	t.Log("Given the need to test the status endpoint.")
	{
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			t.Fatal("\tShould be able to create a request.",
				ballotX, err)
		}
		t.Log("\tShould be able to create a request.",
			checkMark)

		rw := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rw, req)

		if rw.Code != 200 {
			t.Fatal("\tShould receive \"200\"", ballotX, rw.Code)
		}
		t.Log("\tShould receive \"200\"", checkMark)
	}
}

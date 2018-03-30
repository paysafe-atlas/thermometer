package main

import (
	"net/http"
	"encoding/json"
	"fmt"
)

var history []string

type testStruct struct {
	Temperature string
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		handleGet(rw, request)
	} else if request.Method == "POST" {
		handlePost(rw, request)
	} else {
		rw.WriteHeader(405)
	}
}

func handlePost(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var t testStruct
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	// f, err := strconv.ParseFloat("3.1415", 64)
	fmt.Println(t.Temperature)
	history = append(history, t.Temperature)
}

func handleGet(rw http.ResponseWriter, request *http.Request) {
	var encoder = json.NewEncoder(rw)
	var err = encoder.Encode(history)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/temperature/log", parseGhPost)
	http.HandleFunc("/temperature/log", parseGhPost)
	http.ListenAndServe(":80", nil)
}

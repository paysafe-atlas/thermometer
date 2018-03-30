package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

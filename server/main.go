package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var history []temperatureLog

type temperatureLog struct {
	Temperature string `json:"temperature"`
	DateCreated time.Time `json:"dateCreated"`
}

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

	var tempLog temperatureLog
	tempLog.Temperature = t.Temperature
	tempLog.DateCreated = time.Now()
	history = append(history, tempLog)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

func handleGet(rw http.ResponseWriter, request *http.Request) {
	var b, err = json.Marshal(history)
	if err != nil {
		panic(err)
	}
	rw.Write(b)
}

func main() {
	http.HandleFunc("/temperature/log", parseGhPost)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

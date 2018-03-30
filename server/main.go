package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"
	"encoding/csv"
	"bufio"
	"io"
)

var history []temperatureLog
var historyFile *os.File
const fileLocation = "historyFile.csv"

type temperatureLog struct {
	Temperature string `json:"temperature"`
	DateCreated time.Time `json:"dateCreated"`
}

type testStruct struct {
	Temperature string
}

func main() {
	loadHistory()
	http.HandleFunc("/temperature/log", parseGhPost)
	http.HandleFunc("/temperature/log/last", getLastLogEntry)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded file and server started on port 8080!");
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		handleGet(rw, request)
	} else if request.Method == "POST" {
		handlePost(rw, request)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
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
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.Write([]byte("{}"))
}

func handleGet(rw http.ResponseWriter, request *http.Request) {
	var b, err = json.Marshal(history)
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.Write(b)
	rw.WriteHeader(http.StatusOK)
}

func getLastLogEntry(rw http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var lastEntry = history[len(history) - 1]
	var b, err = json.Marshal(lastEntry)
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.Write(b)
	rw.WriteHeader(http.StatusOK)
}

func loadHistory() {
	historyFile, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE, 600)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(bufio.NewReader(historyFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			panic(error)
		}
		var dateCreated, _ = time.Parse(time.RFC3339, line[1])
		history = append(history, temperatureLog{
			Temperature: line[0],
			DateCreated: dateCreated,
		})
	}
}

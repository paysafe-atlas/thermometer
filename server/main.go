package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var history []temperatureLog

const fileLocation = "historyFile.csv"

type temperatureLog struct {
	Temperature string    `json:"temperature"`
	DateCreated time.Time `json:"dateCreated"`
}

type testStruct struct {
	Temperature string
}

func main() {
	loadHistory()
	http.HandleFunc("/temperature/log", getAllEntries)
	http.HandleFunc("/temperature/log/last", getLastLogEntry)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded file and server started on port 8080!")
}

func getAllEntries(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
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

	fmt.Println(t.Temperature)

	var tempLog temperatureLog
	tempLog.Temperature = t.Temperature
	tempLog.DateCreated = time.Now()
	history = append(history, tempLog)
	logEntryInFile(tempLog)

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.Write([]byte("{}"))
}

func handleGet(rw http.ResponseWriter, request *http.Request) {
	var queryParams = request.URL.Query()
	var limitValues = queryParams["limit"]
	var historyLen = len(history)
	var limit = historyLen

	if limitValues != nil {
		limit64, _ := strconv.ParseInt(limitValues[0], 10, 64)
		limit = int(limit64)
		if limit > historyLen {
			limit = historyLen
		}
	}

	var b, err = json.Marshal(history[historyLen-limit : historyLen])
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}

func getLastLogEntry(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if request.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var lastEntry = history[len(history)-1]
	var b, err = json.Marshal(lastEntry)
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	rw.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	rw.Header().Set("pragma", "no-cache")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}

func loadHistory() {
	var historyFile = getHistoryFile()
	reader := csv.NewReader(bufio.NewReader(historyFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			panic(error)
		}
		var dateCreated, _ = time.Parse(time.RFC3339Nano, line[1])
		history = append(history, temperatureLog{
			Temperature: line[0],
			DateCreated: dateCreated,
		})
	}
	historyFile.Close()
}

func logEntryInFile(entry temperatureLog) {
	var historyFile = getHistoryFile()
	var line = fmt.Sprintf("\n%v,%v", entry.Temperature, entry.DateCreated.Format(time.RFC3339Nano))
	var _, error = historyFile.WriteString(line)
	historyFile.Close()
	if error != nil {
		log.Fatal(error)
	}
}

func getHistoryFile() *os.File {
	var historyFile, err = os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		panic(err)
	}
	return historyFile
}

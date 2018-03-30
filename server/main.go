package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type testStruct struct {
	Temperature string
}

func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var t testStruct
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	// f, err := strconv.ParseFloat("3.1415", 64)
	fmt.Println(t.Temperature)
}

func main() {
	http.HandleFunc("/", parseGhPost)
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func MockHandler(data string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, data)
	}
}

func main() {
	jobsData, _ := ioutil.ReadFile("jobs.json")
	http.HandleFunc("/jobs", MockHandler(string(jobsData)))
	log.Println("Starting server on localhost:8081")
	http.ListenAndServe(":8081", nil)
}

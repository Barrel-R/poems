package main

import (
	"fmt"
	"log"
	"net/http"
)

type Poem struct {
	title      string
	path       string
	created_at uint64
}

func createServer() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "%s\n", "Hello from crab")
	})

	log.Println("Server created at localhost:9090")
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("Error while creating the server: ", err)
	}
}

func main() {
	createServer()
}

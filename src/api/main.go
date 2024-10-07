package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Poem struct {
	title      string
	path       string
	created_at uint64
}

func GetPoems(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting poems")
}

func ShowPoem(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "showing poem")
}

func CreatePoem(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "creating poem")
}

func EditPoem(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "editing poem")
}

func DeletePoem(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "deleting poem")
}

func createServer() {
	r := mux.NewRouter()

	r.HandleFunc("/poems", GetPoems).Methods("GET")
	r.HandleFunc("/poems", CreatePoem).Methods("POST")
	r.HandleFunc("/poems/{id}", ShowPoem).Methods("GET")
	r.HandleFunc("/poems/{id}", EditPoem).Methods("PUT")
	r.HandleFunc("/poems/{id}", DeletePoem).Methods("DELETE")

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal("Error while creating server:", err)
	}

	fmt.Println("Server running at localhost:8080")
}

func main() {
	createServer()
}

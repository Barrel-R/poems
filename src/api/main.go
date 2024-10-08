package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Poem struct {
	Id         uint32 `json:"id"`
	Title      string `json:"title"`
	Path       string `json:"path"` // TODO : change path to load contents of file
	Created_at string `json:"created_at"`
}

type ApiResponse struct {
	Data    []Poem `json:"data"`
	Message string `json:"message"`
	Status  uint   `json:"status"`
}

func GetPoems(w http.ResponseWriter, r *http.Request) {
	poemsFile, err := os.ReadFile("../storage/poemas.json")

	if err != nil {
		log.Fatal("Error while getting poems: ", err)
	}

	var poems []Poem

	if err := json.Unmarshal(poemsFile, &poems); err != nil {
		log.Fatal("Couldn't marshal the poems: ", err)
	}

	res := ApiResponse{poems, "List of poems retrieved successfully", http.StatusOK}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RawPoem struct {
	Id         uint32 `json:"id"`
	Title      string `json:"title"`
	Path       string `json:"path"`
	Created_at string `json:"created_at"`
}

type Poem struct {
	Id         uint32    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
}

type ApiResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  uint        `json:"status"`
}

type NotFoundError struct {
	arg     uint32
	message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.message)
}

func getPoemContent(rawPoem RawPoem) (Poem, error) {
	content, err := os.ReadFile("../storage/" + rawPoem.Path)

	if err != nil {
		return Poem{}, fmt.Errorf("couldn't read the poem text file: %w", err)
	}

	date, err := time.Parse(time.DateOnly, rawPoem.Created_at)

	if err != nil {
		return Poem{}, fmt.Errorf("Error while parsing the poem date: %w", err)
	}

	poem := Poem{rawPoem.Id, rawPoem.Title, string(content), date}

	return poem, err
}

func FindPoem(poemId uint32, poems []RawPoem) (RawPoem, error) {
	for _, poem := range poems {
		if poem.Id == poemId {
			return poem, nil
		}
	}

	return RawPoem{}, &NotFoundError{poemId, "couldn't find poem with id "}
}

func ReadPoemFile(pRawPoems *[]RawPoem) {
	poemsFile, err := os.ReadFile("../storage/poemas.json")

	if err != nil {
		log.Fatal("Couldn't read the poems file: ", err)
	}

	if err := json.Unmarshal(poemsFile, pRawPoems); err != nil {
		log.Fatal("Couldn't marshal the poems: ", err)
	}
}

func LoadPoemContents(poems *[]Poem, rawPoems []RawPoem) {
	for _, rawPoem := range rawPoems {
		poem, err := getPoemContent(rawPoem)

		if err != nil {
			log.Fatal("Couldn't get the poem content: ", err)
		}

		*poems = append(*poems, poem)
	}
}

func GetPoems(w http.ResponseWriter, r *http.Request) {
	var rawPoems []RawPoem
	var poems []Poem

	ReadPoemFile(&rawPoems)
	LoadPoemContents(&poems, rawPoems)

	res := ApiResponse{poems, "List of poems retrieved successfully", http.StatusOK}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

func ShowPoem(w http.ResponseWriter, r *http.Request) {
	var rawPoems []RawPoem
	vars := mux.Vars(r)
	poemId, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.Fatal("Couldn't parse the Poem ID: ", err)
	}

	ReadPoemFile(&rawPoems)

	rawPoem, err := FindPoem(uint32(poemId), rawPoems)
	poem, err := getPoemContent(rawPoem)

	res := ApiResponse{poem, "Poem retrieved successfully", http.StatusOK}

	if err != nil {
		log.Fatal("Couldn't get the poem content: ", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
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

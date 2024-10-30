package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

func getGreatestId() uint32 {
	var poems []RawPoem

	ReadPoemFile(&poems)

	id := uint32(0)

	for _, poem := range poems {
		if poem.Id > id {
			id = poem.Id
		}
	}

	return id
}

func createPoemTextFile(poem Poem) error {
	data, err := json.Marshal(poem.Content)

	if err != nil {
		fmt.Printf("error while marshaling the poem: %v\n", err)
		return err
	}

	err = os.WriteFile("../storage/"+strings.Trim(strings.ToLower(poem.Title), " ")+".txt", data, 0644)

	return nil
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
		log.Fatalf("Couldn't read the poems file: %v\n", err)
	}

	if err := json.Unmarshal(poemsFile, pRawPoems); err != nil {
		log.Fatalf("Couldn't marshal the poems: %v\n", err)
	}
}

func LoadPoemContents(poems *[]Poem, rawPoems []RawPoem) {
	for _, rawPoem := range rawPoems {
		poem, err := getPoemContent(rawPoem)

		if err != nil {
			log.Fatalf("Couldn't get the poem content: %v\n", err)
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
		log.Fatalf("Couldn't parse the Poem ID: %v\n", err)
	}

	ReadPoemFile(&rawPoems)

	rawPoem, err := FindPoem(uint32(poemId), rawPoems)
	poem, err := getPoemContent(rawPoem)

	res := ApiResponse{poem, "Poem retrieved successfully", http.StatusOK}

	if err != nil {
		log.Fatalf("Couldn't get the poem content: %v\n", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

func CreatePoem(w http.ResponseWriter, r *http.Request) {
	id := getGreatestId() + 1
	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) == 0 {
		log.Fatal("The title must not be empty.\n")
	}

	if len(content) == 0 {
		log.Fatal("The content must not be empty.\n")
	}

	poem := Poem{id, title, content, time.Now()}

	err := createPoemTextFile(poem)

	if err != nil {
		log.Printf("error while saving poem text file: %v\n", err)
		return
	}

	path := strings.Trim(strings.ToLower(poem.Title), " ") + ".txt"

	rawPoem := RawPoem{id, poem.Title, path, time.Now().Format("2006-01-02")}

	var jsonData []RawPoem
	ReadPoemFile(&jsonData)

	jsonData = append(jsonData, rawPoem)

	data, err := json.Marshal(jsonData)

	if err != nil {
		log.Printf("error while marshaling the new created poem: %v\n", err)
		return
	}

	err = os.WriteFile("../storage/poemas.json", data, 0644)

	if err != nil {
		log.Printf("error while saving the new JSON: %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ApiResponse{poem, "Poem created successfully", http.StatusOK}

	json.NewEncoder(w).Encode(res)
}

func EditPoem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	parsedId, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		log.Fatal("Couldn't parse the poem ID\n")
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(title) == 0 {
		log.Fatal("The title must not be empty.\n")
	}

	if len(content) == 0 {
		log.Fatal("The content must not be empty.\n")
	}

	poem := Poem{uint32(parsedId), title, content, time.Now()}

	err = createPoemTextFile(poem)

	if err != nil {
		log.Printf("error while saving poem text file: %v\n", err)
		return
	}

	path := strings.Trim(strings.ToLower(poem.Title), " ") + ".txt"

	rawPoem := RawPoem{uint32(parsedId), title, path, time.Now().Format("2006-01-02")}

	var jsonData []RawPoem
	ReadPoemFile(&jsonData)

	for i, poem := range jsonData {
		if poem.Id == uint32(parsedId) {
			jsonData[i] = rawPoem
		}
	}

	data, err := json.Marshal(jsonData)

	if err != nil {
		log.Printf("error while marshaling the new created poem: %v\n", err)
		return
	}

	err = os.WriteFile("../storage/poemas.json", data, 0644)

	if err != nil {
		log.Printf("error while saving the new JSON: %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ApiResponse{poem, "Poem edited successfully", http.StatusOK}

	json.NewEncoder(w).Encode(res)
}

func DeletePoem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	parsedId, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		log.Fatal("Couldn't parse the poem ID\n")
	}

	var jsonData []RawPoem
	ReadPoemFile(&jsonData)

	for i, poem := range jsonData {
		if poem.Id == uint32(parsedId) {
			jsonData = slices.Delete(jsonData, i, i+1)
			// err := os.Remove("../storage/" + poem.Path)

			// NOTE: add some sort of confirmation to remove the file?

			// if err != nil {
			// 	log.Printf("Couldn't remove the poem file: %v\n", err)
			// }
		}
	}

	data, err := json.Marshal(jsonData)

	if err != nil {
		log.Printf("error while marshaling the data: %v\n", err)
		return
	}

	err = os.WriteFile("../storage/poemas.json", data, 0644)

	if err != nil {
		log.Printf("error while saving the new JSON: %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ApiResponse{nil, "Poem deleted successfully", http.StatusOK}

	json.NewEncoder(w).Encode(res)
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

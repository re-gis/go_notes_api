package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	Id      int    `json:"Id"`
	Content string `json:"content"`
}

var notes = []Note{}
var IdCounter int = 0

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var t Note
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Content == "" {
		fmt.Println(t.Content)
	}

	IdCounter++
	t.Id = IdCounter
	notes = append(notes, t)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(notes)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var updatedNote Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for index, t := range notes {
		if t.Id == updatedNote.Id {
			notes[index] = updatedNote
			json.NewEncoder(w).Encode(updatedNote)
			return
		}
	}

	http.Error(w, "Task not found!", http.StatusNotFound)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	IdStr := parts[len(parts)-1]
	fmt.Println(IdStr)
	if IdStr == "" {
		http.Error(w, "Id parameter is required!", http.StatusBadRequest)
		return
	}
	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for index, t := range notes {
		if t.Id == Id {
			notes = append(notes[:index], notes[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Task not found!", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/note", getNotesHandler)
	http.HandleFunc("/note/create", createNoteHandler)
	http.HandleFunc("/note/update", updateNoteHandler)
	http.HandleFunc("/note/delete/", deleteNote)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

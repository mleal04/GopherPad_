package handlers

import (
	"encoding/json"
	"net/http"
	"notes-api/models"
	"notes-api/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GetAllNotes --> GET
func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := storage.AllNotes()
	//write back the http response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notes)
}

// CreateNote --> POST
func CreateNote(w http.ResponseWriter, r *http.Request) {
	//create a note struct
	var note models.Note
	//decode http request body into the note struct
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//generate a new UUID for the note
	note.ID = uuid.New().String()
	//store the note in memory --> sending a full correct json struct
	storage.Create(note)
	//write back the http response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetNote --> GET
func GetNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note, err := storage.GetNoteByID(id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	//write back the http response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// UpdateNote --> PUT
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var note models.Note
	//decode http request body into the note struct
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	note.ID = id
	//send this upddate note to the storage
	if err := storage.UpdateNote(note); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	//write back the http response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// DeleteNote --> Delete
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := storage.DeleteNoteById(id); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	//write back the http response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}

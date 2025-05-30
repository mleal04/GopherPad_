package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes-api/models"
	"notes-api/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GetAllNotes --> GET
func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GET request for all notes")
	notes := storage.AllNotes()
	w.WriteHeader(http.StatusOK) // 200 OK for GET
	json.NewEncoder(w).Encode(notes)
}

// CreateNote --> POST
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	note.ID = uuid.New().String()
	storage.Create(note)
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
	w.WriteHeader(http.StatusOK) // 200 OK for GET
	json.NewEncoder(w).Encode(note)
}

// UpdateNote --> PUT
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	note.ID = id
	if err := storage.UpdateNote(note); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// DeleteNote --> DELETE
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := storage.DeleteNoteById(id); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}

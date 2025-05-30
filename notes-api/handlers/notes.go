package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes-api/auth"
	"notes-api/models" // models: Note and LoginRequest
	"notes-api/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// curl -X POST http://localhost:8080/login \
//      -H "Content-Type: application/json" \
//      -d '{"username": "admin", "password": "password"}'

// curl -X POST http://localhost:8080/notes \
//   -H "Authorization: Bearer $TOKEN" \
//   -H "Content-Type: application/json" \
//   -d '{"title": "Test Note", "content": "This is a secured note"}'

// curl http://localhost:8080/notes \
//   -H "Authorization: Bearer $TOKEN"

// curl http://localhost:8080/notes/<id> \
//   -H "Authorization: Bearer $TOKEN"

// curl -X PUT http://localhost:8080/notes/<id> \
//   -H "Authorization: Bearer $TOKEN" \
//   -H "Content-Type: application/json" \
//   -d '{"title": "Updated Title", "content": "Updated content"}'

// Login --> POST to verify the user and return a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	//get user credentials from the request body
	log.Println("Received POST request for login")
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//verify the credentials
	if req.Username != "admin" || req.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	//generate a token (in this case, just a simple string)
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	//send the token back to the user
	w.WriteHeader(http.StatusOK) // 200 OK for successful login
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetAllNotes --> GET
func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GET request for all notes")
	notes := storage.AllNotes()
	w.WriteHeader(http.StatusOK) // 200 OK for GET
	json.NewEncoder(w).Encode(notes)
}

// CreateNote --> POST
func CreateNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Received POST request to create a new note")
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
	log.Println("Received GET request for a specific note")
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
	log.Println("Received PUT request to update a note")
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
	log.Println("Received DELETE request to delete a note")
	id := chi.URLParam(r, "id")
	if err := storage.DeleteNoteById(id); err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}

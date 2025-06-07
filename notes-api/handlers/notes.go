package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes-api/auth"
	"notes-api/models"  // models
	"notes-api/storage" //DB1 and DB2
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// curl -X POST http://localhost:8080/register \
//      -H "Content-Type: application/json" \
//      -d '{"username": "glkee", "password": "hellodad"}'

// curl -X POST http://localhost:8080/login \
//      -H "Content-Type: application/json" \
//      -d '{"username": "glkee", "password": "hellodad"}'

// curl -X POST http://localhost:8080/notes/glkee \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imdsa2VlIiwiZXhwIjoxNzQ5MzIyNjIyfQ.7MoTIWVlBz_oYb_eEbkxUj8dmx7X6r7a6mTp0jBNzzk" \
//   -H "Content-Type: application/json" \
//   -d '{"username": "glkee", "title": "my-dad-phone", "content": "340-555-1234"}'

// curl http://localhost:8080/notes/glkee \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imdsa2VlIiwiZXhwIjoxNzQ5MzIyNjIyfQ.7MoTIWVlBz_oYb_eEbkxUj8dmx7X6r7a6mTp0jBNzzk"

// curl http://localhost:8080/notes/glkee/2 \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imdsa2VlIiwiZXhwIjoxNzQ5MzIyNjIyfQ.7MoTIWVlBz_oYb_eEbkxUj8dmx7X6r7a6mTp0jBNzzk"

// curl -X POST http://localhost:8080/notes/glkee/2 \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imdsa2VlIiwiZXhwIjoxNzQ5MzIyNjIyfQ.7MoTIWVlBz_oYb_eEbkxUj8dmx7X6r7a6mTp0jBNzzk"
//   -d '{"title": "my-dad-phone", "content": "4444444444"}'

// curl -X DELETE http://localhost:8080/notes/mleal2/<id> \
//   -H "Authorization: Bearer $TOKEN" \
//   -H "Content-Type: application/json" \

// CreateUser --> POST to register a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	//get user credentials from the request body
	log.Println("Received POST request for user registration")
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	//create a new user
	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}
	if err := storage.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered"))
	log.Println("User registered successfully")
}

// Login --> POST to verify the user and return a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	//get user credentials from the request body
	log.Println("Received POST request for login")
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var user models.User //should be filled with the correct user data
	if err := storage.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
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
	log.Println("User logged in successfully, token generated")
}

// GetAllNotes --> GET
func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GET request for all notes")
	//get the username
	username := chi.URLParam(r, "username")
	log.Printf("Retrieving notes for user: %s", username)
	//check the username exists in the DB1
	var user models.User
	if err := storage.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	log.Printf("User %s found, retrieving notes", user.Username)
	//get all notes from DB2
	var notes []models.Notes
	if err := storage.DB2.Where("username = ?", user.Username).Find(&notes).Error; err != nil {
		log.Printf("Error retrieving notes for user %s: %v", user.Username, err)
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
	log.Println("Notes retrieved successfully from DB2")
}

// CreateNote --> POST
func CreateNote(w http.ResponseWriter, r *http.Request) {
	//token should be checked before this handler is called
	log.Println("Received POST request to create a new note")
	//get the username, title and content from the request body
	var input models.NoteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//create the note
	note := models.Notes{
		Username: input.Username,
		Title:    input.Title,
		Content:  input.Content,
	}
	//save the note to DB2
	if err := storage.DB2.Create(&note).Error; err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}
	//send the response back to the user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
	log.Println("Note created successfully to DB2")
}

// GetNote --> GET
func GetNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GET request for a specific note")

	// grab the username and the note-id
	username := chi.URLParam(r, "username")
	idStr := chi.URLParam(r, "id")

	// convert the note-id to an int
	noteIndex, err := strconv.Atoi(idStr)
	if err != nil || noteIndex < 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// check if the user exists in DB1
	var user models.User
	if err := storage.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// get all the notes based on the username --> which we know exists now
	var notes []models.Note
	if err := storage.DB2.Where("username = ?", username).Order("id").Find(&notes).Error; err != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}

	// Bounds check
	if noteIndex >= len(notes) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Get the correct note
	correctNote := notes[noteIndex]

	// Return the note
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(correctNote)
}

// UpdateNote --> PUT
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Received PUT request to update a note")
	//get the username and the note-id
	username := chi.URLParam(r, "username")
	id := chi.URLParam(r, "id")
	noteIndex, err := strconv.Atoi(id)
	if err != nil || noteIndex < 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	//check the username exists in the DB1
	var user models.User
	if err := storage.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	//get all the notes from DB2 and check if the note exists
	var notes []models.Note
	if err := storage.DB2.Where("username = ?", username).Order("id").Find(&notes).Error; err != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}
	//bounds check
	if noteIndex >= len(notes) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	// Get the correct note
	correctNote := notes[noteIndex]
	// update the note
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	correctNote.Title = input.Title
	correctNote.Content = input.Content
	fmt.Println(correctNote.Title)
	// save the note to DB2
	if err := storage.DB2.Save(&correctNote).Error; err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}
	//send the response back to the user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(correctNote)
	log.Println("Note updated successfully to DB2")
}

// DeleteNote --> DELETE
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	log.Println("Received DELETE request to delete a note")
	username := chi.URLParam(r, "username")
	id := chi.URLParam(r, "id")
	noteIndex, err := strconv.Atoi(id)
	if err != nil || noteIndex < 0 {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	//check the username exists in the DB1
	var user models.User
	if err := storage.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	//get all the notes from DB2 and check if the note exists
	var notes []models.Note
	if err := storage.DB2.Where("username = ?", username).Order("id").Find(&notes).Error; err != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}
	//bounds check
	if noteIndex >= len(notes) {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	// Get the correct note
	correctNote := notes[noteIndex]
	// delete the note
	if err := storage.DB2.Delete(&correctNote).Error; err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}
	//send the response back to the user
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
	log.Println("Note deleted successfully from DB2")
}

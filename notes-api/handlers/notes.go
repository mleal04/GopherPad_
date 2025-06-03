package handlers

import (
	"encoding/json"
	"net/http"
	"notes-api/auth"
	"notes-api/models"  // models
	"notes-api/storage" //DB1 and DB2

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// curl -X POST http://localhost:8080/register \
//      -H "Content-Type: application/json" \
//      -d '{"username": "admin", "password": "trial"}'

// curl -X POST http://localhost:8080/login \
//      -H "Content-Type: application/json" \
//      -d '{"username": "mleal2", "password": "hellomom"}'

// curl -X POST http://localhost:8080/notes/mleal2 \
//   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1sZWFsMiIsImV4cCI6MTc0ODkzNTU3Nn0.BiM86cC_-yLVaohDJe0bNWS1m0J9pbc6TKtOrWnmTFM" \
//   -H "Content-Type: application/json" \
//   -d '{"username": "mleal2", "title": "my-first-note-2", "content": "this is my first note-2"}'

// curl http://localhost:8080/notes/mleal2 \                                                             ─╯
//  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1sZWFsMiIsImV4cCI6MTc0ODkzNTU3Nn0.BiM86cC_-yLVaohDJe0bNWS1m0J9pbc6TKtOrWnmTFM" \

// curl http://localhost:8080/notes/mleal2/<id> \
//   -H "Authorization: Bearer $TOKEN"
// 	-d '{"usernane": mlea2}'

// curl -X PUT http://localhost:8080/notes/mleal2/<id> \
//   -H "Authorization: Bearer $TOKEN" \
//   -H "Content-Type: application/json" \
//   -d '{"title": "Updated Title", "content": "Updated content"}'

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

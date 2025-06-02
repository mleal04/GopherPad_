package main

import (
	"log"
	"net/http"
	"notes-api/auth"
	"notes-api/handlers"
	"notes-api/models"  // models: Note and LoginRequest
	"notes-api/storage" //use connect to connect to the database

	"github.com/go-chi/chi/v5"
)

// curl -X POST http://localhost:8080/register \
//      -H "Content-Type: application/json" \
//      -d '{"username": "admin", "password": "trial"}'

// curl -X POST http://localhost:8080/login \
//      -H "Content-Type: application/json" \
//      -d '{"username": "admin", "password": "trial"}'

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

func main() {
	//connect to the database (first)
	if err := storage.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("DB connected successfully 🎉")

	//set up the database schema --> if it doesn't exist
	if err := storage.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate DB schema:", err)
	}
	log.Println("DB schema migrated successfully 🎉")

	//create a router for the server
	r := chi.NewRouter()
	//route to login
	r.Post("/register", handlers.CreateUser)
	r.Post("/login", handlers.Login)

	//group the /notes routes together
	r.Group(func(r chi.Router) {
		//always check the JWT token for the /notes routes
		r.Use(auth.JWTMiddleware)

		r.Route("/notes", func(r chi.Router) {
			r.Get("/", handlers.GetAllNotes)
			r.Post("/", handlers.CreateNote)
			r.Get("/{id}", handlers.GetNote)
			r.Put("/{id}", handlers.UpdateNote)
			r.Delete("/{id}", handlers.DeleteNote)
		})
	})
	//start the server --> where r is the router
	http.ListenAndServe(":8080", r)
}

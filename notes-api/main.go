package main

import (
	"net/http"
	"notes-api/auth"
	"notes-api/handlers"

	"github.com/go-chi/chi/v5"
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

func main() {
	//create a router for the server
	r := chi.NewRouter()
	//route to login
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

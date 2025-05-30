package main

import (
	"net/http"
	"notes-api/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	//create a router for the server
	r := chi.NewRouter()
	//create the routes
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", handlers.GetAllNotes)       //get all notes in storage
		r.Post("/", handlers.CreateNote)       //post a new note
		r.Get("/{id}", handlers.GetNote)       //get a specific note
		r.Put("/{id}", handlers.UpdateNote)    //update a specific note
		r.Delete("/{id}", handlers.DeleteNote) //delete a note
	})
	//start the server
	http.ListenAndServe(":8080", r)
}

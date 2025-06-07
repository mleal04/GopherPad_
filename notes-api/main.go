package main

import (
	// "log"
	"net/http"
	"notes-api/auth"
	"notes-api/handlers"
	"notes-api/models"  // models: Note and LoginRequest
	"notes-api/storage" //use connect to connect to the database

	"os"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
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

// curl http://localhost:8080/notes/mleal2/3 \                                                             ─╯
//  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1sZWFsMiIsImV4cCI6MTc0ODkzNTU3Nn0.BiM86cC_-yLVaohDJe0bNWS1m0J9pbc6TKtOrWnmTFM" \

// curl http://localhost:8080/notes/mleal2/3 \
//   -H "Authorization: Bearer $TOKEN"

// curl -X PUT http://localhost:8080/notes/mleal2/3 \
//   -H "Authorization: Bearer $TOKEN" \
//   -H "Content-Type: application/json" \
//   -d '{"title": "Updated Title", "content": "Updated content"}'

// start the logger
func logging_entry() {
	// Create or append to logs/app.log
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Log in JSON for Splunk-friendly format
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func main() {
	//initialize the logger
	logging_entry()
	log.WithFields(log.Fields{
		"username": "admin",
		"event":    "login_success",
	}).Info("User logged in")

	//connect to the databases (first)
	if err := storage.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("DB's connected successfully 🎉")

	//set up the database schema --> if it doesn't exist
	if err := storage.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate DB schema:", err)
	}
	log.Println("DB1 schema migrated successfully 🎉")
	if err := storage.DB2.AutoMigrate(&models.Notes{}); err != nil {
		log.Fatal("Failed to migrate DB schema:", err)
	}
	log.Println("DB2 schema migrated successfully 🎉")

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
			r.Get("/{username}", handlers.GetAllNotes)
			r.Post("/{username}", handlers.CreateNote)
			r.Get("/{username}/{id}", handlers.GetNote)    //this will get the note (mleal2/3)
			r.Put("/{username}/{id}", handlers.UpdateNote) //this will update the note (mleal2/3)
			r.Delete("/{username}/{id}", handlers.DeleteNote)
		})
	})
	//start the server --> where r is the router
	http.ListenAndServe(":8080", r)
}

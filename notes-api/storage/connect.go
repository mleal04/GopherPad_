package storage

import (
	"log"

	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// variable DB will hold the database connection
var DB *gorm.DB

// connect to the existing database using GORM
func Connect() error {
	// Load environment variables from .env file
	err2 := godotenv.Load()
	if err2 != nil {
		log.Println("Error loading .env file, using default connection string")
	}
	dsn := os.Getenv("DATABASE_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to DB:", err)
	}
	return err
}

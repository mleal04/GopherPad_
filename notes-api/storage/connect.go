package storage

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// variable DB will hold the database connection
var DB *gorm.DB
var DB2 *gorm.DB

// connect to the existing database using GORM
func Connect() error {
	// Load environment variables from .env file
	err2 := godotenv.Load()
	if err2 != nil {
		log.Println("Error loading .env file, using default connection string")
	}
	//connect to DB1
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to DB:", err)
	}
	log.Println("connected to mydb successfully")
	//connect to DB2
	dsn2 := os.Getenv("DATABASE_URL2")
	var err3 error
	DB2, err3 = gorm.Open(postgres.Open(dsn2), &gorm.Config{})
	if err3 != nil {
		log.Println("Failed to connect to DB:", err3)
	}
	log.Println("connected to mydb2 successfully")
	return err
}

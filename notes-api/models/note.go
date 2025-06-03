package models

import "github.com/golang-jwt/jwt/v5"

// schema for the database --> user login
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string // this will store the hashed password
}

// schema for the database --> notes
type Notes struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"index"` // foreign key-like, but just indexed
	Title    string
	Content  string
}

type NoteCreateRequest struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// schema for the notes on the database
type Note struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// use this for gateway struct
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// this will be the model for the JWT token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

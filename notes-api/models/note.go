package models

import "github.com/golang-jwt/jwt/v5"

// schema for the database
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string // this will store the hashed password
}

// this will be the model for the database
type Note struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// this will be the model for login authentication
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// this will be the model for the JWT token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

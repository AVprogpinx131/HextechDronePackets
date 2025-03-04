package api

import (
	"encoding/json"
	"hextech_interview_project/internal/auth"
	"hextech_interview_project/internal/models"
	"hextech_interview_project/internal/repository"
	"log"
	"net/http"
	"database/sql"
)

// Register a new user
func RegisterHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err = repository.RegisterUser(db, creds.Username, creds.Password)
    if err != nil {
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login and get JWT token
func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        log.Println("Invalid request payload:", err)
        return
    }

    userID, err := repository.AuthenticateUser(db, creds.Username, creds.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        log.Println("Invalid credentials for user:", creds.Username, "-", err)
        return
    }

    token, err := auth.GenerateJWT(userID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        log.Println("Error generating token:", err)
        return
    }

    log.Println("User logged in successfully:", creds.Username)

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

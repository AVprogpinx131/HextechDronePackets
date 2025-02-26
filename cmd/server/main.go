package main

import (
	"hextech_interview_project/internal/api"
	"hextech_interview_project/internal/auth"
	"hextech_interview_project/internal/repository"
	"hextech_interview_project/internal/websocket"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
    // Initialize PostgreSQL database
    repository.InitDB()

    r := mux.NewRouter() 

    // Public routes (no JWT token required)
    r.HandleFunc("/register", api.RegisterHandler).Methods("POST")
    r.HandleFunc("/login", api.LoginHandler).Methods("POST")

    // Protected routes (require JWT token)
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(auth.JWTMiddleware)
    protected.HandleFunc("/protected", api.ProtectedHandler).Methods("GET")

    // Set up WebSocket handler
    r.HandleFunc("/ws", websocket.HandleWebSocket)

    // Start HTTP server
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

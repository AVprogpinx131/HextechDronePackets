package main

import (
	"hextech_interview_project/internal/repository"
	"hextech_interview_project/internal/websocket"
	"log"
	"net/http"
)

func main() {
    // Initialize PostgreSQL database
    repository.InitDB()

    // Set up WebSocket handler
    http.HandleFunc("/ws", websocket.HandleWebSocket)

    // Start HTTP server
    log.Println("Starting server on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

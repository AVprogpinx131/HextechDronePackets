package websocket

import (
    "github.com/gorilla/websocket"
    "io"
    "log"
    "net/http"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// HandleWebSocket upgrades an HTTP request to a WebSocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    RegisterClient(conn) // Register client in the hub
    log.Println("New WebSocket client connected")

    // Listen for messages
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            if err == io.EOF {
                log.Println("Client disconnected")
            } else {
                log.Println("Error reading message:", err)
            }
            break
        }

        NotifyUsers("Received message: " + string(message)) // Broadcast message

        ProcessPacket(message) // Process the received packet
    }

    UnregisterClient(conn) // Remove client when disconnected
}
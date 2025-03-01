package websocket

import (
    "sync"
    "github.com/gorilla/websocket"
)

// WebSocket Hub (tracks connected users)
type Hub struct {
    clients map[int]*websocket.Conn // Maps userID → WebSocket connection
    mutex   sync.Mutex
}

var hub = Hub{
    clients: make(map[int]*websocket.Conn),
}

// Register a new WebSocket client (user)
func RegisterClient(userID int, conn *websocket.Conn) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()
    hub.clients[userID] = conn
}

// Remove a WebSocket client when they disconnect
func UnregisterClient(userID int) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()
    delete(hub.clients, userID)
}

// Send a message to a specific user
func NotifyUser(userID int, message string) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()

    conn, exists := hub.clients[userID]
    if !exists {
        return
    }

    err := conn.WriteMessage(websocket.TextMessage, []byte(message))
    if err != nil {
        conn.Close()
        delete(hub.clients, userID)
    }
}

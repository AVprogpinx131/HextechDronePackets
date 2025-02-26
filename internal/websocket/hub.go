package websocket

import (
    "sync"
    "github.com/gorilla/websocket"
)

// Hub struct manages all WebSocket clients and messages
type Hub struct {
    clients map[*websocket.Conn]bool // Active WebSocket connections
    mutex   sync.Mutex               // Ensures thread safety
}

// Create a global hub instance
var hub = Hub{
    clients: make(map[*websocket.Conn]bool),
}

// RegisterClient adds a new WebSocket connection to the hub
func RegisterClient(conn *websocket.Conn) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()
    hub.clients[conn] = true
}

// UnregisterClient removes a WebSocket connection from the hub
func UnregisterClient(conn *websocket.Conn) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()
    delete(hub.clients, conn)
}

// BroadcastMessage sends a message to all connected clients
func BroadcastMessage(msg []byte) {
    hub.mutex.Lock()
    defer hub.mutex.Unlock()

    for conn := range hub.clients {
        err := conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            conn.Close()
            delete(hub.clients, conn)
        }
    }
}

// NotifyUsers sends a message to all connected users
func NotifyUsers(message string) {
    BroadcastMessage([]byte(message))
}

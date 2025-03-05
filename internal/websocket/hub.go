package websocket

import (
    "log"
    "sync"
    "github.com/gorilla/websocket"
)

var (
    userConnections = make(map[int]*websocket.Conn)
    mutex           sync.Mutex
)


// Register a new WebSocket client
func RegisterClient(userID int, conn *websocket.Conn) {
    mutex.Lock()
    defer mutex.Unlock()
    userConnections[userID] = conn
    log.Printf("WebSocket client registered: User %d", userID)
}


// Remove a WebSocket client when they disconnect
func UnregisterClient(userID int) {
    mutex.Lock()
    defer mutex.Unlock()
    delete(userConnections, userID)
    log.Printf("WebSocket client unregistered: User %d", userID)
}


// Send a message to a specific user
func NotifyUser(userID int, message []byte) {
    mutex.Lock()
    conn, exists := userConnections[userID]
    mutex.Unlock()

    if exists {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Printf("Error notifying user %d: %v", userID, err)
            UnregisterClient(userID)
            conn.Close()
        }
    }
}
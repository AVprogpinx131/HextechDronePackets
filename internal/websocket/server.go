package websocket

import (
    "hextech_interview_project/internal/auth"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// Handle WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    log.Println("New WebSocket connection request")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade WebSocket:", err)
        http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
        return
    }

    // Extract JWT token from query parameters
    token := r.URL.Query().Get("token")
    if token == "" {
        log.Println("Missing token in WebSocket request")
        conn.Close()
        return
    }

    log.Println("Token received:", token)

    // Validate the token and get the user ID
    userID, err := auth.ValidateToken(token)
    if err != nil {
        log.Println("Invalid token:", err)
        conn.Close()
        return
    }

    log.Println("User authenticated with ID:", userID)

    RegisterClient(userID, conn)

    log.Println("WebSocket connection established for user:", userID)

    // Keep WebSocket connection open
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("🔌 WebSocket disconnected: User %d", userID)
            UnregisterClient(userID)
            return
        }
    }

}

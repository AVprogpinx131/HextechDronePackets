package websocket

import (
	"database/sql"
	"fmt"
	"hextech_interview_project/internal/repository"
	"log"
	"sync"
	"time"
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
}


// Send a message to a specific user
func NotifyUser(userID int, message string) {
    mutex.Lock()
    conn, exists := userConnections[userID]
    mutex.Unlock()

    if exists {
        conn.WriteMessage(websocket.TextMessage, []byte(message))
    }
}


// Periodic updates to send active drone information
func StartPeriodicUpdates(db *sql.DB) {
    go func() {
        for {
            time.Sleep(10 * time.Second)

            log.Println("Running periodic drone updates...")

            mutex.Lock()
            log.Printf("Active WebSocket clients: %d", len(userConnections))

            for userID, conn := range userConnections {
                activeDrones, err := repository.GetDronesInsideTerritory(db, userID)
                if err != nil {
                    log.Println("Error fetching drones:", err)
                    continue
                }

                if len(activeDrones) > 0 {
                    log.Printf("Sending %d active drones to user %d", len(activeDrones), userID) 

                    // Format drones as table for logging
                    log.Println("Active Drones:")
                    log.Println("+-------------------+-----------+-----------+---------+")
                    log.Println("| MAC Address      | Latitude  | Longitude | Altitude|")
                    log.Println("+-------------------+-----------+-----------+---------+")
                    for _, d := range activeDrones {
                        log.Printf("| %-17s | %-9.5f | %-9.5f | %-7.2f |", d.MAC, d.Latitude, d.Longitude, d.Altitude)
                    }
                    log.Println("+-------------------+-----------+-----------+---------+")

                    // Send formatted message to WebSocket
                    message := fmt.Sprintf("Active Drones (%d):\n", len(activeDrones))
                    for _, d := range activeDrones {
                        message += fmt.Sprintf("%s at (%.5f, %.5f) Alt: %.2f\n", d.MAC, d.Latitude, d.Longitude, d.Altitude)
                    }

                    err := conn.WriteMessage(websocket.TextMessage, []byte(message))
                    if err != nil {
                        log.Println("Error sending message:", err)
                        conn.Close()
                        delete(userConnections, userID)
                    }
                }
            }
            mutex.Unlock()
        }
    }()
}

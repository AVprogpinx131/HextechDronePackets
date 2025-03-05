package websocket

import (
    "database/sql"
    "encoding/json"
    "hextech_interview_project/internal/auth"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
    "fmt"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// Handle WebSocket connections and drone notifications
func HandleWebSocket(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    log.Println("New WebSocket connection request")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade WebSocket:", err)
        http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
        return
    }

    token := r.URL.Query().Get("token")
    if token == "" {
        log.Println("Missing token in WebSocket request")
        conn.Close()
        return
    }

    userID, err := auth.ValidateToken(token)
    if err != nil {
        log.Println("Invalid token:", err)
        conn.Close()
        return
    }

    log.Println("User authenticated with ID:", userID)

    RegisterClient(userID, conn)
    defer func() {
        UnregisterClient(userID)
        conn.Close()
        log.Printf("WebSocket disconnected: User %d", userID)
    }()

    log.Println("WebSocket connection established for user:", userID)

    go StartPeriodicUpdates(db, userID)

    // Monitor for new drone entries
    knownDrones := make(map[string]struct{})
    entryTicker := time.NewTicker(1 * time.Second)
    defer entryTicker.Stop()

    for range entryTicker.C {
        drones, err := repository.GetDronesInsideTerritory(db, userID)
        if err != nil {
            log.Println("Error fetching drones:", err)
            continue
        }

        for _, drone := range drones {
            key := fmt.Sprintf("%s-%f-%f-%f", drone.MAC, drone.Latitude, drone.Longitude, drone.Altitude)
            if _, exists := knownDrones[key]; !exists {
                knownDrones[key] = struct{}{}
                territories, err := repository.GetTerritories(db, userID)
                if err != nil {
                    log.Println("Error fetching territories:", err)
                    continue
                }
                for _, t := range territories {
                    if IsDroneInsideTerritory(models.DronePacket{
                        MAC:       drone.MAC,
                        Latitude:  drone.Latitude,
                        Longitude: drone.Longitude,
                        Altitude:  drone.Altitude,
                    }, t) {
                        direction := calculateBearing(t.Latitude, t.Longitude, drone.Latitude, drone.Longitude)
                        notify := struct {
                            Event     string  `json:"event"`
                            MAC       string  `json:"mac"`
                            Territory string  `json:"territory"`
                            Direction float64 `json:"direction"`
                        }{
                            Event:     "drone_entered",
                            MAC:       drone.MAC,
                            Territory: t.Name,
                            Direction: direction,
                        }
                        msg, err := json.Marshal(notify)
                        if err != nil {
                            log.Println("Error marshaling entry notification:", err)
                            continue
                        }
                        NotifyUser(userID, msg)
                        break
                    }
                }
            }
        }
    }
}

// Send active drones every 10 seconds
func StartPeriodicUpdates(db *sql.DB, userID int) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        activeDrones, err := repository.GetDronesInsideTerritory(db, userID)
        if err != nil {
            log.Println("Error fetching active drones:", err)
            continue
        }

        if len(activeDrones) > 0 {
            log.Printf("Sending %d active drones to user %d", len(activeDrones), userID)
            log.Println("Active Drones:")
            log.Println("+-------------------+-----------+-----------+---------+")
            log.Println("| MAC Address      | Latitude  | Longitude | Altitude|")
            log.Println("+-------------------+-----------+-----------+---------+")
            for _, d := range activeDrones {
                log.Printf("| %-17s | %-9.5f | %-9.5f | %-7.2f |", d.MAC, d.Latitude, d.Longitude, d.Altitude)
            }
            log.Println("+-------------------+-----------+-----------+---------+")

            type ActiveDronesMessage struct {
                Event  string                   `json:"event"`
                Drones []models.DronePacket     `json:"drones"`
            }

            var dronePackets []models.DronePacket
            for _, d := range activeDrones {
                dronePackets = append(dronePackets, models.DronePacket{
                    MAC:       d.MAC,
                    Latitude:  d.Latitude,
                    Longitude: d.Longitude,
                    Altitude:  d.Altitude,
                })
            }

            message := ActiveDronesMessage{
                Event:  "active_drones",
                Drones: dronePackets,
            }
            msgBytes, err := json.Marshal(message)
            if err != nil {
                log.Println("Error marshaling active drones:", err)
                continue
            }
            NotifyUser(userID, msgBytes)
        }
    }
}
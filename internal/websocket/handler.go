package websocket

import (
    "encoding/json"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "net/http"
    "fmt"
    "database/sql"
)


var activeDrones = make(map[string][]int)


// Handles incoming drone data from WebSocket
func ProcessPacket(db *sql.DB, data []byte) {
    var packet models.DronePacket
    err := json.Unmarshal(data, &packet)
    if err != nil {
        log.Println("Invalid packet:", err)
        return
    }

    territories, err := repository.GetAllTerritories(db)
    if err != nil {
        log.Println("Error fetching territories:", err)
        return
    }

    insideTerritories := []int{}
    for _, territory := range territories {
        if IsDroneInsideTerritory(packet, territory) {
            insideTerritories = append(insideTerritories, territory.ID)
        }
    }

    previousTerritories, exists := activeDrones[packet.MAC]

    if exists {
        // Check for exits
        for _, prevTerritory := range previousTerritories {
            if !contains(insideTerritories, prevTerritory) {
                log.Printf("Drone %s EXITED territory %d", packet.MAC, prevTerritory)
                repository.SaveDroneMovement(db, packet.MAC, prevTerritory, "exit")
                repository.SaveExitEvent(db, packet.MAC)
            }
        }
        
        // Check for new entries
        for _, newTerritory := range insideTerritories {
            if !contains(previousTerritories, newTerritory) {
                log.Printf("Drone %s ENTERED territory %d", packet.MAC, newTerritory)
                repository.SaveDroneMovement(db, packet.MAC, newTerritory, "entry")

                // Notify the territory owner
                ownerID, err := repository.GetTerritoryOwner(db, newTerritory)
                if err == nil {
                    message := fmt.Sprintf("Drone %s entered your territory %d", packet.MAC, newTerritory)
                    NotifyUser(ownerID, message)
                }
            }
        }
    }

    activeDrones[packet.MAC] = insideTerritories
    repository.SavePacket(db, packet)
}


// Check if a value exists in a slice
func contains(slice []int, value int) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}


// Allows HTTP clients to send drone data
func HandleDronePacket(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var packet models.DronePacket

    // Decode JSON body
    err := json.NewDecoder(r.Body).Decode(&packet)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Reuse ProcessPacket logic
    jsonData, _ := json.Marshal(packet)
    ProcessPacket(db, jsonData)

    // Respond with success
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Drone packet received"})
}
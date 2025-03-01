package websocket

import (
    "encoding/json"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "net/http"
)


var activeDrones = make(map[string]int)


// Handles incoming drone data from WebSocket
func ProcessPacket(data []byte) {
    var packet models.DronePacket
    err := json.Unmarshal(data, &packet)
    if err != nil {
        log.Println("Invalid packet:", err)
        return
    }

    territories, err := repository.GetAllTerritories()
    if err != nil {
        log.Println("Error fetching territories:", err)
        return
    }

    insideTerritoryID := 0
    for _, territory := range territories {
        if IsDroneInsideTerritory(packet, territory) {
            insideTerritoryID = territory.ID
            log.Printf("Drone %s is INSIDE territory %s (ID: %d)", packet.MAC, territory.Name, territory.ID)
            break
        }
    }

    previousTerritoryID, exists := activeDrones[packet.MAC]

    if !exists {
        log.Printf("First detection of drone %s - setting initial state.", packet.MAC)
        activeDrones[packet.MAC] = insideTerritoryID
        if insideTerritoryID != 0 {
            log.Printf("First detected INSIDE territory %d - recording ENTRY", insideTerritoryID)
            repository.SaveDroneMovement(packet.MAC, insideTerritoryID, "entry")
        }
    } else {
        if previousTerritoryID != 0 && insideTerritoryID == 0 {
            log.Printf("Drone %s EXITED territory %d", packet.MAC, previousTerritoryID)
            repository.SaveDroneMovement(packet.MAC, previousTerritoryID, "exit")
            repository.SaveExitEvent(packet.MAC) 
        } else if previousTerritoryID == 0 && insideTerritoryID != 0 {
            log.Printf("Drone %s ENTERED territory %d", packet.MAC, insideTerritoryID)
            repository.SaveDroneMovement(packet.MAC, insideTerritoryID, "entry")
        }
    }

    activeDrones[packet.MAC] = insideTerritoryID
    repository.SavePacket(packet)
}


// Allows HTTP clients to send drone data
func HandleDronePacket(w http.ResponseWriter, r *http.Request) {
    var packet models.DronePacket

    // Decode JSON body
    err := json.NewDecoder(r.Body).Decode(&packet)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Reuse ProcessPacket logic
    jsonData, _ := json.Marshal(packet)
    ProcessPacket(jsonData)

    // Respond with success
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Drone packet received"})
}
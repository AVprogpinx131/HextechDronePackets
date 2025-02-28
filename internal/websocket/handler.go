package websocket

import (
    "encoding/json"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "net/http"
)


var activeDrones = make(map[string]bool)


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

    log.Printf("Found %d territories in DB", len(territories))

    insideAny := false
    for _, territory := range territories {
        log.Printf("Checking territory: %s (Lat: %f, Lon: %f, Radius: %f)", 
            territory.Name, territory.Latitude, territory.Longitude, territory.Radius)
        
        if IsDroneInsideTerritory(packet, territory) {
            insideAny = true
            log.Printf("Drone %s is INSIDE territory: %s", packet.MAC, territory.Name)
            break
        }
    }

    if wasInside, exists := activeDrones[packet.MAC]; exists && wasInside && !insideAny {
        log.Printf("Drone %s EXITED a territory!", packet.MAC)
        repository.SaveExitEvent(packet.MAC)
    }

    activeDrones[packet.MAC] = insideAny

    err = repository.SavePacket(packet)
    if err != nil {
        log.Println("Error saving packet:", err)
    }
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
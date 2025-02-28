package websocket

import (
    "encoding/json"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "net/http"
)

// Handles incoming drone data from WebSocket
func ProcessPacket(data []byte) {
    var packet models.DronePacket

    // Decode the JSON payload
    err := json.Unmarshal(data, &packet)
    if err != nil {
        log.Println("Invalid packet:", err)
        return
    }

    // Save packet to the database
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
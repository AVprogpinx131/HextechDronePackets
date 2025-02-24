package websocket

import (
    "encoding/json"
    "log"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
)

func ProcessPacket(data []byte) {
    var packet models.DronePacket
    err := json.Unmarshal(data, &packet)
    if err != nil {
        log.Println("Invalid packet:", err)
        return
    }
    repository.SavePacket(packet) // Save packet in the database
}
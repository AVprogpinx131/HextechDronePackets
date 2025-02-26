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
        log.Println("Invalid packet format:", err)
        return
    }

    log.Println("Received packet:", packet)

    // Save packet to the database
    if err := repository.SavePacket(packet); err != nil {
        log.Println("Error saving packet:", err)
    } else {
        log.Println("Packet saved successfully:", packet)
    }

    // Notify users about the new drone event
    NotifyUsers("Drone detected: " + packet.MAC)
}
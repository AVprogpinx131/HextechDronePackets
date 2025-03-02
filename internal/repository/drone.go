package repository

import (
    "log"
    "database/sql"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/config"
    "fmt"
)

var db *sql.DB


func InitDB() {
    config.LoadConfig()

    var err error
    db, err = sql.Open("postgres", config.DbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }

    fmt.Println("Successfully connected to PostgreSQL")
}


func SavePacket(packet models.DronePacket) error {
	query := `INSERT INTO drone_packets (mac, latitude, longitude, altitude) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, packet.MAC, packet.Latitude, packet.Longitude, packet.Altitude)
	if err != nil {
        log.Println("Error inserting packet:", err)
    }
    return err
}


func SaveExitEvent(mac string) error {
    query := `INSERT INTO drone_exits (mac, exit_time) VALUES ($1, NOW())`
    _, err := db.Exec(query, mac)
    if err != nil {
        log.Println("Error saving exit event:", err)
    } else {
        log.Printf("Saved exit event: MAC=%s", mac)
    }
    return err
}


func SaveDroneMovement(mac string, territoryId int, eventType string) error {
    query := `INSERT INTO drone_movements (mac, territory_id, event_type) VALUES ($1, $2, $3)`
    _, err := db.Exec(query, mac, territoryId, eventType)
    if err != nil {
        log.Println("Error saving drone movement:", err)
    } else {
        log.Printf("Saved drone movement: MAC=%s, Territory=%d, Event=%s", mac, territoryId, eventType)
    }
    return err
}
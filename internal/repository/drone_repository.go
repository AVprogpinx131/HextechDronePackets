package repository

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/config"
    "fmt"
)

var db *sql.DB


func InitDB() {
    // Load configuration
    co.LoadConfig()

    // Connect to PostgreSQL
    var err error
    db, err = sql.Open("postgres", co.DbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Verify connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }

    fmt.Println("Successfully connected to PostgreSQL")
}

func SavePacket(packet models.DronePacket) error {
	query := `INSERT INTO drone_packets (mac, lat, lon, altitude) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, packet.MAC, packet.Lat, packet.Lon, packet.Altitude)
	if err != nil {
        log.Println("Error inserting packet:", err)
    }
    return err
}

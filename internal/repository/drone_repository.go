package repository

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "os"
    "hextech_interview_project/internal/models"
    "fmt"
    "github.com/joho/godotenv"
)

var db *sql.DB


func InitDB() {
    // Load .env file
    err := godotenv.Load("config/.env")
    if err != nil {
        log.Println("Warning: No .env file found")
    }

    // Read DATABASE_URL
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set in .env file")
    }

    // Connect to PostgreSQL
    db, err = sql.Open("postgres", dbURL)
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
	return err
}

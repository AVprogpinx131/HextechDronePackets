package config

import (
    "os"
    "github.com/joho/godotenv"
    "log"
)

var (
	JwtSecret []byte
	DbURL string
)

func LoadConfig() {
    // Load .env file
    err := godotenv.Load("config/.env")
    if err != nil {
        log.Println("Warning: No .env file found")
    }

    // Read JWT_SECRET
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set in .env file")
    }

    JwtSecret = []byte(secret)

	// Read DATABASE_URL
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set in .env file")
    }
	DbURL = dbURL
}
package main

import (
	"hextech_interview_project/internal/api"
	"hextech_interview_project/internal/auth"
	"hextech_interview_project/internal/repository"
	"hextech_interview_project/internal/websocket"
	"log"
	"net/http"
	"github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

func main() {
    db, err := repository.InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Initialize schema
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS territories (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            name TEXT NOT NULL,
            latitude DOUBLE PRECISION NOT NULL,
            longitude DOUBLE PRECISION NOT NULL,
            radius DOUBLE PRECISION NOT NULL,
            min_altitude DOUBLE PRECISION NOT NULL,
            max_altitude DOUBLE PRECISION NOT NULL
        );
        CREATE TABLE IF NOT EXISTS drone_packets (
            mac TEXT,
            latitude DOUBLE PRECISION,
            longitude DOUBLE PRECISION,
            altitude DOUBLE PRECISION
        );
        CREATE TABLE IF NOT EXISTS drone_exits (
            id SERIAL PRIMARY KEY,
            mac TEXT,
            exit_time TIMESTAMP WITH TIME ZONE
        );
        CREATE TABLE IF NOT EXISTS drone_movements (
            id SERIAL PRIMARY KEY,
            mac TEXT,
            territory_id INTEGER REFERENCES territories(id),
            event_type TEXT,
            timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        );
    `)

    if err != nil {
        log.Fatal("Failed to initialize schema:", err)
    }
    log.Println("Database schema initialized")


    r := mux.NewRouter() 

    // Public routes (no JWT token required)
    r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        api.RegisterHandler(db, w, r)
    }).Methods("POST")
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        api.LoginHandler(db, w, r)
    }).Methods("POST")

    // Protected routes (require JWT token in Authorization header)
    protected := r.PathPrefix("/api").Subrouter()
    protected.Use(auth.JWTMiddleware)
    protected.HandleFunc("/protected", api.ProtectedHandler).Methods("GET")

    // Territory routes
    protected.HandleFunc("/territories", func(w http.ResponseWriter, r *http.Request) {
        api.CreateTerritoryHandler(db, w, r)
    }).Methods("POST")
    protected.HandleFunc("/territories", func(w http.ResponseWriter, r *http.Request) {
        api.GetTerritoriesHandler(db, w, r)
    }).Methods("GET")
    protected.HandleFunc("/territories/{id}", func(w http.ResponseWriter, r *http.Request) {
        api.DeleteTerritoryHandler(db, w, r)
    }).Methods("DELETE")

    // Report routes
    protected.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
        api.GetDroneMovements(db, w, r)
    }).Methods("GET")

    // Drone packet routes
    protected.HandleFunc("/drone_packet", func(w http.ResponseWriter, r *http.Request) {
        websocket.HandleDronePacket(db, w, r)
    }).Methods("POST")

    // Set up WebSocket handler
    r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        websocket.HandleWebSocket(db, w, r)
    })

    // Start HTTP server
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))

    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

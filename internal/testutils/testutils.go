package testutils

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/lib/pq"
)

// Creates a PostgreSQL test container and returns a connected *sql.DB and a cleanup function
func SetupTestDB() (*sql.DB, func()) {
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "postgres",
            "POSTGRES_PASSWORD": "testpassword",
            "POSTGRES_DB":       "hextech_test_db",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(120 * time.Second),
    }
    postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        panic("Failed to start test container: " + err.Error())
    }

    host, _ := postgresC.Host(ctx)
    port, _ := postgresC.MappedPort(ctx, "5432/tcp")
    dbURL := fmt.Sprintf("postgres://postgres:testpassword@%s:%s/hextech_test_db?sslmode=disable", host, port.Port())

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        panic("Failed to connect to test DB: " + err.Error())
    }

    for i := 0; i < 15; i++ {
        if err = db.Ping(); err == nil {
            break
        }
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        panic("Failed to ping test DB after retries: " + err.Error())
    }

    // Setup schema for all tables
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
        panic("Failed to create test tables: " + err.Error())
    }

    cleanup := func() {
        postgresC.Terminate(ctx)
    }
    return db, cleanup
}

// Helper function to clear all tables in the test database and create a test user

func ClearTables(db *sql.DB) error {
    _, err := db.Exec("TRUNCATE TABLE users, territories, drone_packets, drone_exits, drone_movements RESTART IDENTITY CASCADE")
    return err
}

func CreateTestUser(db *sql.DB, username, password string) (int, error) {
    var userID int
    err := db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, password).Scan(&userID)
    if err != nil {
        return 0, fmt.Errorf("failed to create test user: %v", err)
    }
    return userID, nil
}
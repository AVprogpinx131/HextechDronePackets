package repository

import (
    "hextech_interview_project/internal/models"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestSavePacket(t *testing.T) {
    t.Run("Valid packet save", func(t *testing.T) {
        packet := models.DronePacket{
            MAC:       "00:11:22:33:44:55",
            Latitude:  40.7128,
            Longitude: -74.0060,
            Altitude:  50.0,
        }
        err := SavePacket(testDB, packet)
        assert.NoError(t, err)

        var count int
        err = testDB.QueryRow("SELECT COUNT(*) FROM drone_packets WHERE mac = $1", packet.MAC).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    })
}

func TestSaveExitEvent(t *testing.T) {
    t.Run("Valid exit event save", func(t *testing.T) {
        mac := "00:11:22:33:44:66"
        err := SaveExitEvent(testDB, mac)
        assert.NoError(t, err)

        var exitTime time.Time
        err = testDB.QueryRow("SELECT exit_time FROM drone_exits WHERE mac = $1", mac).Scan(&exitTime)
        assert.NoError(t, err)
        assert.WithinDuration(t, time.Now(), exitTime, 5*time.Second)
    })
}

func TestSaveDroneMovement(t *testing.T) {
    t.Run("Valid movement save", func(t *testing.T) {
        // Setup: Create a user
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 1, "testuser1", "hashedpass")
        assert.NoError(t, err)

        // Setup: Create a territory
        territory := models.Territory{
            UserID:      1,
            Name:        "Test Territory",
            Latitude:    40.0,
            Longitude:   40.0,
            Radius:      1000.0,
            MinAltitude: 0.0,
            MaxAltitude: 100.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var territoryID int
        err = testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Test Territory").Scan(&territoryID)
        assert.NoError(t, err)

        // Save movement
        mac := "00:11:22:33:44:77"
        eventType := "ENTER"
        err = SaveDroneMovement(testDB, mac, territoryID, eventType)
        assert.NoError(t, err)

        var savedEventType string
        err = testDB.QueryRow("SELECT event_type FROM drone_movements WHERE mac = $1 AND territory_id = $2", mac, territoryID).Scan(&savedEventType)
        assert.NoError(t, err)
        assert.Equal(t, eventType, savedEventType)
    })
}
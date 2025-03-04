package repository

import (
    "hextech_interview_project/internal/models"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestGetMovementsByUser(t *testing.T) {
    t.Run("Retrieve user movements", func(t *testing.T) {
        // Setup: Create a user
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 2, "testuser2", "hashedpass")
        assert.NoError(t, err)

        // Setup: Create a territory
        territory := models.Territory{
            UserID:      2,
            Name:        "Movement Territory",
            Latitude:    50.0,
            Longitude:   50.0,
            Radius:      1000.0,
            MinAltitude: 0.0,
            MaxAltitude: 100.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var territoryID int
        err = testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Movement Territory").Scan(&territoryID)
        assert.NoError(t, err)

        // Insert movements
        err = SaveDroneMovement(testDB, "00:11:22:33:44:88", territoryID, "ENTER")
        assert.NoError(t, err)
        time.Sleep(1 * time.Second)
        err = SaveDroneMovement(testDB, "00:11:22:33:44:88", territoryID, "EXIT")
        assert.NoError(t, err)

        movements, err := GetMovementsByUser(testDB, 2)
        assert.NoError(t, err)
        assert.Len(t, movements, 2)

        assert.Equal(t, "EXIT", movements[0].EventType)
        assert.Equal(t, "ENTER", movements[1].EventType)
        assert.Equal(t, "Movement Territory", movements[0].TerritoryName)
        assert.Equal(t, territoryID, movements[0].TerritoryId)
        assert.Equal(t, 100.0, movements[0].MaxAltitude)
        assert.Equal(t, 0.0, movements[0].MinAltitude)
    })

    t.Run("No movements for user", func(t *testing.T) {
        movements, err := GetMovementsByUser(testDB, 999)
        assert.NoError(t, err)
        assert.Len(t, movements, 0)
    })
}
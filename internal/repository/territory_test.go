package repository

import (
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/testutils"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCreateTerritory(t *testing.T) {
    t.Run("Valid territory creation", func(t *testing.T) {
        assert.NoError(t, testutils.ClearTables(testDB)) // Clear tables
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 1, "testuser1", "hashedpass")
        assert.NoError(t, err)

        territory := models.Territory{
            UserID:      1,
            Name:        "Test Territory",
            Latitude:    40.7128,
            Longitude:   -74.0060,
            Radius:      1000.0,
            MinAltitude: 0.0,
            MaxAltitude: 500.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var name string
        err = testDB.QueryRow("SELECT name FROM territories WHERE user_id = $1 AND name = $2", 1, "Test Territory").Scan(&name)
        assert.NoError(t, err)
        assert.Equal(t, "Test Territory", name)
    })
}

func TestGetDronesInsideTerritory(t *testing.T) {
    t.Run("Drones within territory", func(t *testing.T) {
        assert.NoError(t, testutils.ClearTables(testDB)) // Clear tables
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 3, "testuser3", "hashedpass")
        assert.NoError(t, err)

        territory := models.Territory{
            UserID:      3,
            Name:        "Drone Territory",
            Latitude:    50.0,
            Longitude:   50.0,
            Radius:      1000.0,
            MinAltitude: 0.0,
            MaxAltitude: 100.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var territoryID int
        err = testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Drone Territory").Scan(&territoryID)
        assert.NoError(t, err)

        _, err = testDB.Exec("INSERT INTO drone_packets (mac, latitude, longitude, altitude) VALUES ($1, $2, $3, $4)",
            "00:11:22:33:44:55", 50.0, 50.0, 50.0)
        assert.NoError(t, err)
        _, err = testDB.Exec("INSERT INTO drone_packets (mac, latitude, longitude, altitude) VALUES ($1, $2, $3, $4)",
            "00:11:22:33:44:66", 60.0, 60.0, 50.0)
        assert.NoError(t, err)

        drones, err := GetDronesInsideTerritory(testDB, 3)
        assert.NoError(t, err)
        assert.Len(t, drones, 1)
        assert.Equal(t, "00:11:22:33:44:55", drones[0].MAC)
    })
}

func TestGetTerritoryOwner(t *testing.T) {
    t.Run("Valid territory owner", func(t *testing.T) {
        assert.NoError(t, testutils.ClearTables(testDB)) // Clear tables
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 4, "testuser4", "hashedpass")
        assert.NoError(t, err)

        territory := models.Territory{
            UserID:      4,
            Name:        "Owner Test",
            Latitude:    10.0,
            Longitude:   10.0,
            Radius:      100.0,
            MinAltitude: 0.0,
            MaxAltitude: 100.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var territoryID int
        err = testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Owner Test").Scan(&territoryID)
        assert.NoError(t, err)

        userID, err := GetTerritoryOwner(testDB, territoryID)
        assert.NoError(t, err)
        assert.Equal(t, 4, userID)
    })
}

func TestDeleteTerritory(t *testing.T) {
    t.Run("Delete existing territory", func(t *testing.T) {
        assert.NoError(t, testutils.ClearTables(testDB)) // Clear tables
        _, err := testDB.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 5, "testuser5", "hashedpass")
        assert.NoError(t, err)

        territory := models.Territory{
            UserID:      5,
            Name:        "Delete Test",
            Latitude:    20.0,
            Longitude:   20.0,
            Radius:      100.0,
            MinAltitude: 0.0,
            MaxAltitude: 100.0,
        }
        err = CreateTerritory(testDB, territory)
        assert.NoError(t, err)

        var territoryID int
        err = testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Delete Test").Scan(&territoryID)
        assert.NoError(t, err)

        err = DeleteTerritory(testDB, 5, territoryID)
        assert.NoError(t, err)

        var count int
        err = testDB.QueryRow("SELECT COUNT(*) FROM territories WHERE id = $1", territoryID).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 0, count)
    })
}
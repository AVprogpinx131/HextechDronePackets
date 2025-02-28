package repository

import (
    "hextech_interview_project/internal/models"
    "log"
    "errors"
    "database/sql"
)

// Create a new territory
func CreateTerritory(userID int, name string, lat, lon, radius float64) error {
    query := `INSERT INTO territories (user_id, name, latitude, longitude, radius) VALUES ($1, $2, $3, $4, $5)`
    _, err := db.Exec(query, userID, name, lat, lon, radius)
    if err != nil {
        log.Println("Error creating territory:", err)
    } else {
        log.Println("Territory created for user:", userID)
    }
    return err
}

// Get all territories for a user
func GetTerritories(userID int) ([]models.Territory, error) {
    query := `SELECT * FROM territories WHERE user_id = $1`
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var territories []models.Territory
    for rows.Next() {
        var t models.Territory
        err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Latitude, &t.Longitude, &t.Radius)
        if err != nil {
            return nil, err
        }
        territories = append(territories, t)
    }
    return territories, nil
}


// Retrieves all territories in the database
func GetAllTerritories() ([]models.Territory, error) {
    query := `SELECT * FROM territories`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var territories []models.Territory
    for rows.Next() {
        var t models.Territory
        err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Latitude, &t.Longitude, &t.Radius)
        if err != nil {
            return nil, err
        }
        territories = append(territories, t)
    }
    return territories, nil
}


// Delete a territory by ID
func DeleteTerritory(userID, territoryID int) error {
    query := `DELETE FROM territories WHERE id = $1 AND user_id = $2 RETURNING id`
    var deletedID int
    err := db.QueryRow(query, territoryID, userID).Scan(&deletedID)

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No territory found with ID: %d for user: %d", territoryID, userID)
            return errors.New("territory not found")
        }
        log.Println("❌ Error deleting territory:", err)
        return err
    }

    log.Printf("Territory deleted: %d", deletedID)
    return nil
}

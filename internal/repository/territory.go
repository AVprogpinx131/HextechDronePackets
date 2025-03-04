package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"hextech_interview_project/internal/models"
	"log"
)

// Create a new territory
func CreateTerritory(db *sql.DB, territory models.Territory) error {
	query := `INSERT INTO territories (user_id, name, latitude, longitude, radius, min_altitude, max_altitude) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query, territory.UserID, territory.Name, territory.Latitude, territory.Longitude, territory.Radius, territory.MinAltitude, territory.MaxAltitude)
	if err != nil {
		log.Println("Error creating territory:", err)
	} else {
		log.Printf("Territory created: %s (Lat: %.6f, Lon: %.6f, Radius: %.2f, Alt: %.2f - %.2f)",
			territory.Name, territory.Latitude, territory.Longitude, territory.Radius,
			territory.MinAltitude, territory.MaxAltitude)
	}
	return err
}

// Get all territories for a user
func GetTerritories(db *sql.DB, userID int) ([]models.Territory, error) {
	query := `SELECT * FROM territories WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var territories []models.Territory
	for rows.Next() {
		var t models.Territory
		err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Latitude, &t.Longitude, &t.Radius, &t.MinAltitude, &t.MaxAltitude)
		if err != nil {
			return nil, err
		}
		territories = append(territories, t)
	}
	return territories, nil
}

// Retrieves all territories in the database
func GetAllTerritories(db *sql.DB) ([]models.Territory, error) {
	query := `SELECT * FROM territories`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var territories []models.Territory
	for rows.Next() {
		var t models.Territory
		err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Latitude, &t.Longitude, &t.Radius, &t.MinAltitude, &t.MaxAltitude)
		if err != nil {
			return nil, err
		}
		territories = append(territories, t)
	}
	return territories, nil
}

// Get all distinct drone packets currently inside a user's territories
func GetDronesInsideTerritory(db *sql.DB, userID int) ([]models.DronePacket, error) {
	query := `
        SELECT DISTINCT ON (dp.mac, dp.latitude, dp.longitude, dp.altitude) 
            dp.mac, dp.latitude, dp.longitude, dp.altitude
        FROM drone_packets dp
        JOIN territories t ON (
            dp.latitude BETWEEN t.latitude - (t.radius / 111111) AND t.latitude + (t.radius / 111111)
            AND dp.longitude BETWEEN t.longitude - (t.radius / (111111 * COS(RADIANS(t.latitude)))) 
            AND dp.longitude + (t.radius / (111111 * COS(RADIANS(t.latitude))))
            AND dp.altitude BETWEEN t.min_altitude AND t.max_altitude
        )
        WHERE t.user_id = $1
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uniqueDrones := make(map[string]models.DronePacket)
	for rows.Next() {
		var drone models.DronePacket
		err := rows.Scan(&drone.MAC, &drone.Latitude, &drone.Longitude, &drone.Altitude)
		if err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s-%f-%f-%f", drone.MAC, drone.Latitude, drone.Longitude, drone.Altitude)
		uniqueDrones[key] = drone
	}

	var drones []models.DronePacket
	for _, drone := range uniqueDrones {
		drones = append(drones, drone)
	}
	return drones, nil
}

// Get a territory by ID
func GetTerritoryOwner(db *sql.DB, territoryId int) (int, error) {
	query := `SELECT user_id FROM territories WHERE id = $1`
	var userId int
	err := db.QueryRow(query, territoryId).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// Delete a territory by ID
func DeleteTerritory(db *sql.DB, userID, territoryID int) error {
	query := `DELETE FROM territories WHERE id = $1 AND user_id = $2 RETURNING id`
	var deletedID int
	err := db.QueryRow(query, territoryID, userID).Scan(&deletedID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No territory found with ID: %d for user: %d", territoryID, userID)
			return errors.New("territory not found")
		}
		log.Println("Error deleting territory:", err)
		return err
	}

	log.Printf("Territory deleted: %d", deletedID)
	return nil
}

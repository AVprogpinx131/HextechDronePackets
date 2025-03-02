package repository

import (
    "hextech_interview_project/internal/models"
    "log"
)

func GetMovementsByUser(userID int) ([]models.DroneMovement, error) {
    query := `
    SELECT dm.id, dm.mac, dm.event_type, dm.territory_id, t.name AS territory_name, dm.timestamp, t.min_altitude, t.max_altitude
    FROM drone_movements dm
    JOIN territories t ON dm.territory_id = t.id
    WHERE t.user_id = $1
    ORDER BY dm.timestamp DESC
    `
    rows, err := db.Query(query, userID)
    if err != nil {
        log.Println("Error fetching drone movements:", err)
        return nil, err
    }
    defer rows.Close()

    var movements []models.DroneMovement
    for rows.Next() {
        var m models.DroneMovement
        err := rows.Scan(&m.ID, &m.MAC, &m.EventType, &m.TerritoryId, &m.TerritoryName, &m.Timestamp, &m.MinAltitude, &m.MaxAltitude)
        if err != nil {
            log.Println("Error scanning drone movement:", err)
            continue
        }
        movements = append(movements, m)
    }

    log.Printf("Returning %d drone movements", len(movements))
    return movements, nil
}

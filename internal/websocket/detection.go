package websocket

import (
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "math"
    "log"
    "database/sql"
)


func IsDroneInsideTerritory(drone models.DronePacket, territory models.Territory) bool {
    const EarthRadius = 6371000 

    dLat := (drone.Latitude - territory.Latitude) * (math.Pi / 180)
    dLon := (drone.Longitude - territory.Longitude) * (math.Pi / 180)

    lat1 := territory.Latitude * (math.Pi / 180)
    lat2 := drone.Latitude * (math.Pi / 180)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    distance := EarthRadius * c 

    altitudeInside := (drone.Altitude >= territory.MinAltitude) && (drone.Altitude <= territory.MaxAltitude)

    log.Printf("Distance from drone %s to territory %s: %.2f meters", drone.MAC, territory.Name, distance)
    log.Printf("Territory radius: %.2f meters", territory.Radius)
    log.Printf("Altitude check: drone altitude=%.2f, min=%.2f, max=%.2f, inside=%t",
        drone.Altitude, territory.MinAltitude, territory.MaxAltitude, altitudeInside)

    if distance <= territory.Radius && altitudeInside {
        log.Printf("Drone %s is INSIDE territory %s", drone.MAC, territory.Name)
        return true
    } else {
        log.Printf("Drone %s is OUTSIDE territory %s", drone.MAC, territory.Name)
        return false
    }
}


func CheckDroneInTerritories(db *sql.DB, drone models.DronePacket) []int {
    territories, err := repository.GetAllTerritories(db)
    if err != nil {
        return nil
    }

    affectedUsers := []int{}
    for _, territory := range territories {
        if IsDroneInsideTerritory(drone, territory) {
            affectedUsers = append(affectedUsers, territory.UserID)
        }
    }
    return affectedUsers
}

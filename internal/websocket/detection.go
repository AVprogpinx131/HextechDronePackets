package websocket

import (
    "database/sql"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "log"
    "math"
)


func IsDroneInsideTerritory(drone models.DronePacket, territory models.Territory) bool {
    const EarthRadius = 6371000

    // Convert degrees to radians
    dLat := (drone.Latitude - territory.Latitude) * (math.Pi / 180)
    dLon := (drone.Longitude - territory.Longitude) * (math.Pi / 180)
    lat1 := territory.Latitude * (math.Pi / 180)
    lat2 := drone.Latitude * (math.Pi / 180)

    // Haversine formula
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    distance := EarthRadius * c

    altitudeInside := drone.Altitude >= territory.MinAltitude && drone.Altitude <= territory.MaxAltitude

    log.Printf("Distance from drone %s to territory %s: %.2f meters", drone.MAC, territory.Name, distance)
    log.Printf("Territory radius: %.2f meters", territory.Radius)
    log.Printf("Altitude check: drone altitude=%.2f, min=%.2f, max=%.2f, inside=%t",
        drone.Altitude, territory.MinAltitude, territory.MaxAltitude, altitudeInside)

    if distance <= territory.Radius && altitudeInside {
        direction := calculateBearing(territory.Latitude, territory.Longitude, drone.Latitude, drone.Longitude)
        log.Printf("Drone %s is INSIDE territory %s, direction: %.2f° from North", drone.MAC, territory.Name, direction)
        return true
    }
    log.Printf("Drone %s is OUTSIDE territory %s", drone.MAC, territory.Name)
    return false
}


func CheckDroneInTerritories(db *sql.DB, drone models.DronePacket) []int {
    territories, err := repository.GetAllTerritories(db)
    if err != nil {
        log.Println("Error fetching all territories:", err)
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


// Calculate direction relative to North Pole
func calculateBearing(lat1, lon1, lat2, lon2 float64) float64 {
    lat1, lon1 = lat1 * math.Pi / 180, lon1 * math.Pi / 180
    lat2, lon2 = lat2 * math.Pi / 180, lon2 * math.Pi / 180
    dLon := lon2 - lon1
    y := math.Sin(dLon) * math.Cos(lat2)
    x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
    bearing := math.Atan2(y, x)
    return math.Mod(bearing*180/math.Pi+360, 360)
}
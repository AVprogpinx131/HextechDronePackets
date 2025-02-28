package websocket

import (
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "math"
)

// Check if a drone is inside a given territory using Haversine formula
func IsDroneInsideTerritory(drone models.DronePacket, territory models.Territory) bool {
    const EarthRadius = 6371000 // Earth radius in meters

    dLat := (drone.Latitude - territory.Latitude) * (math.Pi / 180)
    dLon := (drone.Longitude - territory.Longitude) * (math.Pi / 180)

    lat1 := territory.Latitude * (math.Pi / 180)
    lat2 := drone.Latitude * (math.Pi / 180)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    distance := EarthRadius * c // Distance in meters

    return distance <= territory.Radius // Check if within radius
}

// Check if a drone is inside any user's territory
func CheckDroneInTerritories(drone models.DronePacket) []int {
    territories, err := repository.GetAllTerritories() // Get all territories
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

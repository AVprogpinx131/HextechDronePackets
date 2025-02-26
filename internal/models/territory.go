package models

type Territory struct {
    ID        int      `json:"id"`
    UserID    int      `json:"user_id"` // Foreign key to users table
    Name      string   `json:"name"`
    Latitude  float64  `json:"latitude"`
    Longitude float64  `json:"longitude"`
    Radius    float64  `json:"radius"` // Defines area size in meters
}

type TerritoryRequest struct {
    Name      string  `json:"name"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Radius    float64 `json:"radius"`
}

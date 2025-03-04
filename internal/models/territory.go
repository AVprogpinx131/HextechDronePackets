package models

type Territory struct {
    ID        int      `json:"id"`
    UserID    int      `json:"user_id"`
    Name      string   `json:"name"`
    Latitude  float64  `json:"latitude"`
    Longitude float64  `json:"longitude"`
    Radius    float64  `json:"radius"`
    MinAltitude float64 `json:"min_altitude"`
    MaxAltitude float64 `json:"max_altitude"`
}

type TerritoryRequest struct {
    Name      string  `json:"name"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Radius    float64 `json:"radius"`
    MinAltitude float64 `json:"min_altitude"`
    MaxAltitude float64 `json:"max_altitude"`
}

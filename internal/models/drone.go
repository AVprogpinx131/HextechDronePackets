package models

type DronePacket struct {
    MAC      string  `json:"mac"`
    Latitude      float64 `json:"latitude"`
    Longitude      float64 `json:"longitude"`
    Altitude float64 `json:"altitude"`
}

type DroneMovement struct {
    ID           int     `json:"id"`
    MAC          string  `json:"mac"`
    EventType    string  `json:"event_type"`
    TerritoryId  int     `json:"territory_id"`
    TerritoryName string `json:"territory_name"`
    Timestamp    string  `json:"timestamp"`
}
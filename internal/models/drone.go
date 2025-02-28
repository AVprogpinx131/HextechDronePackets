package models

type DronePacket struct {
    MAC      string  `json:"mac"`
    Latitude      float64 `json:"latitude"`
    Longitude      float64 `json:"longitude"`
    Altitude float64 `json:"altitude"`
}

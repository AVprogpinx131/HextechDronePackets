package models

type DronePacket struct {
    MAC      string  `json:"mac"`
    Latitude      float64 `json:"lat"`
    Longitude      float64 `json:"lon"`
    Altitude float64 `json:"altitude"`
}

package models

type DronePacket struct {
    MAC      string  `json:"mac"`
    Lat      float64 `json:"lat"`
    Lon      float64 `json:"lon"`
    Altitude float64 `json:"altitude"`
}

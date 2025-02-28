package repository

import (
    "log"
)

func SaveExitEvent(mac string) error {
    query := `INSERT INTO drone_exits (mac, exit_time) VALUES ($1, NOW())`
    _, err := db.Exec(query, mac)
    if err != nil {
        log.Println("Error saving exit event:", err)
    }
    return err
}

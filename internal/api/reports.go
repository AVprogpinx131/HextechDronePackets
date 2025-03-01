package api

import (
    "encoding/json"
    "hextech_interview_project/internal/repository"
	"hextech_interview_project/internal/auth"
    "net/http"
    "log"
)


func GetDroneMovements(w http.ResponseWriter, r *http.Request) {
    userIDRaw := r.Context().Value(auth.UserIDKey) 

    // Ensure userID is correctly retrieved
    userID, ok := userIDRaw.(int)
    if !ok {
        log.Println("Error: user_id is missing or not an integer", userIDRaw)
        http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
        return
    }

    log.Printf("Fetching reports for user_id: %d", userID)

    movements, err := repository.GetMovementsByUser(userID)
    if err != nil {
        log.Println("Failed to fetch reports:", err)
        http.Error(w, "Failed to fetch reports", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movements)
}

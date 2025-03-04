package api

import (
    "encoding/json"
    "hextech_interview_project/internal/auth"
    "net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := auth.GetUserID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "You accessed a protected route!",
        "user_id": userID,
    })
}

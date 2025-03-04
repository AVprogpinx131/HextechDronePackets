package api

import (
    "encoding/json"
    "hextech_interview_project/internal/auth"
    "hextech_interview_project/internal/repository"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
	"hextech_interview_project/internal/models"
    "database/sql"
)

// Create a new territory
func CreateTerritoryHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    userID, err := auth.GetUserID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req models.TerritoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Validate input
    if req.Radius <= 0 || req.MaxAltitude < req.MinAltitude {
        http.Error(w, "Invalid radius or altitude range", http.StatusBadRequest)
        return
    }

    err = repository.CreateTerritory(db, models.Territory{
        UserID:       userID,
        Name:         req.Name,
        Latitude:     req.Latitude,
        Longitude:    req.Longitude,
        Radius:       req.Radius,
        MinAltitude:  req.MinAltitude,
        MaxAltitude:  req.MaxAltitude,
    })

    if err != nil {
        http.Error(w, "Failed to create territory", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Territory created"})
}

// Get all territories for a user
func GetTerritoriesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    userID, err := auth.GetUserID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    territories, err := repository.GetTerritories(db, userID)
    if err != nil {
        http.Error(w, "Failed to fetch territories", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(territories)
}

// Delete a territory
func DeleteTerritoryHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    userID, err := auth.GetUserID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    territoryID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid territory ID", http.StatusBadRequest)
        return
    }

    err = repository.DeleteTerritory(db, userID, territoryID)
    if err != nil {
        http.Error(w, "Failed to delete territory", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Territory deleted"})
}

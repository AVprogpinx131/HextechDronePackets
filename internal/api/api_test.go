package api

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "hextech_interview_project/internal/auth"
    "hextech_interview_project/internal/models"
    "hextech_interview_project/internal/repository"
    "hextech_interview_project/internal/testutils"
    "net/http"
    "net/http/httptest"
    "testing"
    "context"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    var cleanup func()
    testDB, cleanup = testutils.SetupTestDB()
    defer cleanup()
    m.Run()
}

func setupMuxRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        RegisterHandler(testDB, w, r)
    }).Methods("POST")
    router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        LoginHandler(testDB, w, r)
    }).Methods("POST")

    // Protected routes with JWTMiddleware
    protected := router.PathPrefix("").Subrouter()
    protected.Use(auth.JWTMiddleware)
    protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
        ProtectedHandler(w, r)
    }).Methods("GET")
    protected.HandleFunc("/movements", func(w http.ResponseWriter, r *http.Request) {
        GetDroneMovements(testDB, w, r)
    }).Methods("GET")
    protected.HandleFunc("/territories", func(w http.ResponseWriter, r *http.Request) {
        CreateTerritoryHandler(testDB, w, r)
    }).Methods("POST")
    protected.HandleFunc("/territories", func(w http.ResponseWriter, r *http.Request) {
        GetTerritoriesHandler(testDB, w, r)
    }).Methods("GET")
    protected.HandleFunc("/territories/{id}", func(w http.ResponseWriter, r *http.Request) {
        DeleteTerritoryHandler(testDB, w, r)
    }).Methods("DELETE")

    return router
}

func TestRegisterUserHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB)) // Clear tables for isolation

    t.Run("Valid registration", func(t *testing.T) {
        payload := `{"username": "testuser", "password": "securepassword"}`
        req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
        assert.Contains(t, resp.Body.String(), "User registered successfully")
    })

    t.Run("Duplicate user", func(t *testing.T) {
        payload := `{"username": "duplicateuser", "password": "securepassword"}`
        req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
        req1.Header.Set("Content-Type", "application/json")
        resp1 := httptest.NewRecorder()
        router.ServeHTTP(resp1, req1)
        assert.Equal(t, http.StatusCreated, resp1.Code)

        req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
        req2.Header.Set("Content-Type", "application/json")
        resp2 := httptest.NewRecorder()
        router.ServeHTTP(resp2, req2)

        assert.Equal(t, http.StatusConflict, resp2.Code)
        assert.Contains(t, resp2.Body.String(), "User already exists")
    })

    t.Run("Invalid payload", func(t *testing.T) {
        payload := `{"username": "testuser"` // Malformed JSON
        req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "Invalid request payload")
    })
}

func TestLoginHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Valid login", func(t *testing.T) {
        // Register a user first
        err := repository.RegisterUser(testDB, "loginuser", "password123")
        assert.NoError(t, err)

        payload := `{"username": "loginuser", "password": "password123"}`
        req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
        var response map[string]string
        json.Unmarshal(resp.Body.Bytes(), &response)
        assert.Contains(t, response, "token")
        assert.NotEmpty(t, response["token"])
    })

    t.Run("Invalid credentials", func(t *testing.T) {
        payload := `{"username": "loginuser", "password": "wrongpass"}`
        req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusUnauthorized, resp.Code)
        assert.Contains(t, resp.Body.String(), "Invalid credentials")
    })
}

func TestProtectedHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Authorized access", func(t *testing.T) {
        err := repository.RegisterUser(testDB, "protecteduser", "pass")
        assert.NoError(t, err)
        userID, err := repository.AuthenticateUser(testDB, "protecteduser", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)
        t.Logf("Generated JWT: %s", token)

        req, _ := http.NewRequest("GET", "/protected", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        t.Logf("Response Code: %d, Body: %s", resp.Code, resp.Body.String())
        assert.Equal(t, http.StatusOK, resp.Code)
        var response map[string]interface{}
        err = json.Unmarshal(resp.Body.Bytes(), &response)
        assert.NoError(t, err, "Failed to unmarshal response: %s", resp.Body.String())
        assert.Equal(t, "You accessed a protected route!", response["message"])
        assert.Equal(t, float64(userID), response["user_id"])
    })

    t.Run("Unauthorized access", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/protected", nil) // No token

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusUnauthorized, resp.Code)
        assert.Contains(t, resp.Body.String(), "Unauthorized")
    })
}

func TestGetDroneMovements(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Valid movements retrieval", func(t *testing.T) {
        // Setup user and territory
        userID, err := testutils.CreateTestUser(testDB, "movementuser", "pass")
        assert.NoError(t, err)
        territory := models.Territory{UserID: userID, Name: "Test", Latitude: 1.0, Longitude: 1.0, Radius: 100.0, MinAltitude: 0.0, MaxAltitude: 100.0}
        assert.NoError(t, repository.CreateTerritory(testDB, territory))
        var territoryID int
        testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Test").Scan(&territoryID)

        // Add movement
        assert.NoError(t, repository.SaveDroneMovement(testDB, "mac1", territoryID, "ENTER"))

        // Get token
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        req, _ := http.NewRequest("GET", "/movements", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        // Inject userID into context (since GetDroneMovements expects it)
        ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
        req = req.WithContext(ctx)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
        var movements []models.DroneMovement
        json.Unmarshal(resp.Body.Bytes(), &movements)
        assert.Len(t, movements, 1)
        assert.Equal(t, "mac1", movements[0].MAC)
    })

    t.Run("Unauthorized", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/movements", nil) // No token or context

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusUnauthorized, resp.Code)
        assert.Contains(t, resp.Body.String(), "Unauthorized")
    })
}

func TestCreateTerritoryHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Valid creation", func(t *testing.T) {
        userID, err := testutils.CreateTestUser(testDB, "territoryuser", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        payload := `{"name": "New Territory", "latitude": 40.0, "longitude": -74.0, "radius": 1000, "min_altitude": 0, "max_altitude": 500}`
        req, _ := http.NewRequest("POST", "/territories", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
        assert.Contains(t, resp.Body.String(), "Territory created")
    })

    t.Run("Invalid payload", func(t *testing.T) {
        userID, err := testutils.CreateTestUser(testDB, "invaliduser", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        payload := `{"name": "Invalid", "radius": -1}` // Invalid radius
        req, _ := http.NewRequest("POST", "/territories", bytes.NewBuffer([]byte(payload)))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "Invalid radius or altitude range")
    })
}

func TestGetTerritoriesHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Valid retrieval", func(t *testing.T) {
        userID, err := testutils.CreateTestUser(testDB, "getuser", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        territory := models.Territory{UserID: userID, Name: "Get Test", Latitude: 1.0, Longitude: 1.0, Radius: 100.0, MinAltitude: 0.0, MaxAltitude: 100.0}
        assert.NoError(t, repository.CreateTerritory(testDB, territory))

        req, _ := http.NewRequest("GET", "/territories", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
        var territories []models.Territory
        json.Unmarshal(resp.Body.Bytes(), &territories)
        assert.Len(t, territories, 1)
        assert.Equal(t, "Get Test", territories[0].Name)
    })
}

func TestDeleteTerritoryHandler(t *testing.T) {
    router := setupMuxRouter()
    assert.NoError(t, testutils.ClearTables(testDB))

    t.Run("Valid deletion", func(t *testing.T) {
        userID, err := testutils.CreateTestUser(testDB, "deleteuser", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        territory := models.Territory{UserID: userID, Name: "Delete Test", Latitude: 1.0, Longitude: 1.0, Radius: 100.0, MinAltitude: 0.0, MaxAltitude: 100.0}
        assert.NoError(t, repository.CreateTerritory(testDB, territory))
        var territoryID int
        testDB.QueryRow("SELECT id FROM territories WHERE name = $1", "Delete Test").Scan(&territoryID)

        req, _ := http.NewRequest("DELETE", "/territories/"+strconv.Itoa(territoryID), nil)
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
        assert.Contains(t, resp.Body.String(), "Territory deleted")
    })

    t.Run("Invalid territory ID", func(t *testing.T) {
        userID, err := testutils.CreateTestUser(testDB, "deleteinvalid", "pass")
        assert.NoError(t, err)
        token, err := auth.GenerateJWT(userID)
        assert.NoError(t, err)

        req, _ := http.NewRequest("DELETE", "/territories/invalid", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "Invalid territory ID")
    })
}
package repository

import (
    "testing"
)

func TestUserRegistration(t *testing.T) {
    t.Run("Create_and_retrieve_user", func(t *testing.T) {
        if testDB == nil {
            t.Fatal("testDB is nil")
        }
        username := "testuser"
        password := "testpass"
        err := RegisterUser(testDB, username, password)
        if err != nil {
            t.Fatalf("Failed to register user: %v", err)
        }

        // Verify the user was added
        var retrievedUsername string
        err = testDB.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&retrievedUsername)
        if err != nil {
            t.Fatalf("Failed to retrieve user: %v", err)
        }
        if retrievedUsername != username {
            t.Errorf("Expected username '%s', got '%s'", username, retrievedUsername)
        }
    })
}

func TestAuthenticateUser(t *testing.T) {
    t.Run("Valid_credentials", func(t *testing.T) {
        if testDB == nil {
            t.Fatal("testDB is nil")
        }
        username := "authuser"
        password := "authpass"
        err := RegisterUser(testDB, username, password)
        if err != nil {
            t.Fatalf("Failed to register user for auth test: %v", err)
        }

        userID, err := AuthenticateUser(testDB, username, password)
        if err != nil {
            t.Fatalf("Failed to authenticate user: %v", err)
        }
        if userID <= 0 {
            t.Errorf("Expected positive user ID, got %d", userID)
        }
    })

    t.Run("Invalid_password", func(t *testing.T) {
        if testDB == nil {
            t.Fatal("testDB is nil")
        }
        username := "authuser_invalid"
        password := "authpass"
        err := RegisterUser(testDB, username, password)
        if err != nil {
            t.Fatalf("Failed to register user for invalid auth test: %v", err)
        }

        _, err = AuthenticateUser(testDB, username, "wrongpass")
        if err == nil {
            t.Error("Expected authentication to fail with wrong password")
        }
    })
}
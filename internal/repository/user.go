package repository

import (
    "database/sql"
    "errors"
    "log"
    "hextech_interview_project/internal/models"
    "golang.org/x/crypto/bcrypt"
)

// Register a new user
func RegisterUser(db *sql.DB, username, password string) error {
    // Check if the username already exists
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
    err := db.QueryRow(query, username).Scan(&exists)
    if err != nil {
        log.Println("Error checking if user exists:", err)
        return err
    }
    if exists {
        log.Println("User already exists:", username)
        return errors.New("user already exists")
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        return err
    }

    // Insert the new user
    query = `INSERT INTO users (username, password) VALUES ($1, $2)`
    _, err = db.Exec(query, username, string(hashedPassword))
    if err != nil {
        log.Println("Error registering user:", err)
        return err
    }

    log.Println("User registered successfully:", username)
    return nil
}

// Authenticate user and return user ID if successful
func AuthenticateUser(db *sql.DB, username, password string) (int, error) {
    var user models.User

    query := `SELECT id, password FROM users WHERE username=$1`
    row := db.QueryRow(query, username)
    err := row.Scan(&user.ID, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.New("user not found")
        }
        return 0, err
    }

    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return 0, errors.New("invalid password")
    }

    return user.ID, nil
}
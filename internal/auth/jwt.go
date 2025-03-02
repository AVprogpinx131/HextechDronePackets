package auth

import (
    "time"
    "github.com/dgrijalva/jwt-go"
    "hextech_interview_project/config"
)

// Creates a new token for the user
func GenerateJWT(userID int) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
    })

    return token.SignedString(config.JwtSecret)
}
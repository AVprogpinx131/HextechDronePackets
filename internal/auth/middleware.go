package auth

import (
    "context"
    "fmt"
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
	"hextech_interview_project/config"
    "errors"
)

// Context key for storing user ID
type contextKey string

const UserIDKey contextKey = "userID"

// Middleware to verify JWT tokens
func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims := jwt.MapClaims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return config.JwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Extract user ID from token
        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            http.Error(w, "Invalid token payload", http.StatusUnauthorized)
            return
        }

        userID := int(userIDFloat)

        // Store user ID in request context
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Extracts user ID from request context
func GetUserID(r *http.Request) (int, error) {
    userID, ok := r.Context().Value(UserIDKey).(int)
    if !ok {
        return 0, fmt.Errorf("user ID not found in request context")
    }
    return userID, nil
}


// Extracts userID from a JWT token
func ValidateToken(tokenString string) (int, error) {
    claims := jwt.MapClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return config.JwtSecret, nil
    })

    if err != nil || !token.Valid {
        return 0, errors.New("invalid token")
    }

    // Extract user ID from token
    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        return 0, errors.New("invalid token payload")
    }

    return int(userIDFloat), nil
}
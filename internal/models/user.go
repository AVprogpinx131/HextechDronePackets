package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"` // Will be stored as a hashed password
}

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
package main

import (
    "database/sql"
    "encoding/json"
    "net/http"

    _ "github.com/mattn/go-sqlite3"
)

type LoginResponse struct {
    Message string `json:"message"`
}

func login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username := r.FormValue("username")
        password := r.FormValue("password")

        // Check if user exists and password is correct
        var existingUser User
        err := db.QueryRow(`SELECT id, password FROM users WHERE username = ?`, username).Scan(&existingUser.id, &existingUser.password)
        if err != nil {
            http.Error(w, "Invalid username or password", http.StatusBadRequest)
            return
        }

        // Check if password is correct
        if existingUser.password != password {
            http.Error(w, "Invalid username or password", http.StatusBadRequest)
            return
        }

        // Login successful, send JSON response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(LoginResponse{Message: "Login successful"})
    }
}
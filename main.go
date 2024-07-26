package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the database
	db, err := sql.Open("sqlite3", "./subad.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", homeHandler)
	r.Get("/login", loginHandler(db))
	r.Post("/login", loginPostHandler(db))
	r.Get("/register", registerHandler(db))
	r.Post("/register", registerPostHandler(db))
	r.Get("/dashboard", dashboardHandler(db))
	r.Get("/logout", logoutHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Implement home page logic
}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement login page logic
	}
}

func loginPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement login post logic
	}
}

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement register page logic
	}
}

func registerPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement register post logic
	}
}

func dashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement dashboard logic
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic
}

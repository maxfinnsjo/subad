package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"your-project-path/handlers"
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

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Routes
	r.Get("/", h.Home)
	r.Get("/login", h.Login)
	r.Post("/login", h.LoginPost)
	r.Get("/register", h.Register)
	r.Post("/register", h.RegisterPost)
	r.Get("/dashboard", h.Dashboard)
	r.Get("/logout", h.Logout)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

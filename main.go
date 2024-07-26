package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maxfinnsjo/subad/database"
	"github.com/maxfinnsjo/subad/handlers"
	"github.com/maxfinnsjo/subad/sessions"
)

func main() {
	// Initialize the database
	db, err := database.NewDB("./subad.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the session store
	sessionStore := sessions.NewSessionStore()

	// Initialize the router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Initialize handlers
	h := handlers.NewHandler(db, sessionStore)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", h.Home)
		r.Get("/login", h.Login)
		r.Post("/login", h.LoginPost)
		r.Get("/register", h.Register)
		r.Post("/register", h.RegisterPost)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/dashboard", h.Dashboard)
		r.Get("/logout", h.Logout)
		r.Post("/pages", h.CreatePage)
		r.Get("/pages/{id}", h.ViewPage)
		r.Get("/generate-token", h.GenerateToken)
		r.Get("/user-status", h.ViewUserStatus)
		r.Get("/earn-token", h.EarnToken)
		r.Post("/trade-token", h.TradeToken)
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

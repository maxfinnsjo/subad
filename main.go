package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maxfinnsjo/subad/database"
	"github.com/maxfinnsjo/subad/handlers"
	"github.com/maxfinnsjo/subad/sessions"
	"github.com/maxfinnsjo/subad/tokens"
)

func main() {
	log.Println("Starting Subad application...")

	// Initialize the database
	db, err := database.NewDB("./subad.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	// Initialize the database schema
	if err := initDatabase(db); err != nil {
		log.Fatalf("Error initializing database schema: %v", err)
	}
	log.Println("Database schema initialized successfully")

	// Initialize the session store
	sessionStore := sessions.NewSessionStore()
	log.Println("Session store initialized")

	// Initialize the token manager
	tokenManager := tokens.NewTokenManager(db)
	log.Println("Token manager initialized")

	// Initialize the router
	r := chi.NewRouter()
	log.Println("Router initialized")

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	log.Println("Middleware set up")

	// Initialize handlers
	h := handlers.NewHandler(db, sessionStore, tokenManager)
	log.Println("Handlers initialized")

	// Serve static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)
	log.Println("Static file server set up")

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", h.Home)
		r.Get("/login", h.Login)
		r.Post("/login", h.LoginPost)
		r.Get("/register", h.Register)
		r.Post("/register", h.RegisterPost)
	})
	log.Println("Public routes set up")

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
	log.Println("Protected routes set up")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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

func initDatabase(db *database.DB) error {
	schemaSQL, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %v", err)
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema.sql: %v", err)
	}
	return nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

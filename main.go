package main
import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
	"strings"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username := r.FormValue("username")
        password := r.FormValue("password")

        var storedPassword string
        var isAdmin bool
        err := db.QueryRow("SELECT password, is_admin FROM users WHERE username = ?", username).Scan(&storedPassword, &isAdmin)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "Invalid username or password")
            return
        }

        if password != storedPassword {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "Invalid username or password")
            return
        }

        w.WriteHeader(http.StatusOK)
        if isAdmin {
            fmt.Fprintf(w, `Login successful! <a href="/admin">Go to Admin Page</a><script>setTimeout(() => window.location.href = '/admin', 3000);</script>`)
        } else {
            fmt.Fprintf(w, `Login successful! <a href="/sub">Go to Sub Page</a><script>setTimeout(() => window.location.href = '/sub', 3000);</script>`)
        }
    }
}


func registerHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username := r.FormValue("username")
        password := r.FormValue("password")
        isAdmin := r.FormValue("is_admin") == "true"

        if username == "" || password == "" {
            http.Error(w, "Username and password are required", http.StatusBadRequest)
            return
        }

        _, err := db.Exec("INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)", username, password, isAdmin)
        if err != nil {
            if strings.Contains(err.Error(), "UNIQUE constraint failed") {
                http.Error(w, "Username already exists", http.StatusBadRequest)
            } else {
                http.Error(w, "Registration failed", http.StatusInternalServerError)
            }
            return
        }

        w.WriteHeader(http.StatusOK)
        if isAdmin {
            fmt.Fprintf(w, `Registration successful! <a href="/admin">Go to Admin Page</a><script>setTimeout(() => window.location.href = '/admin', 5000);</script>`)
        } else {
            fmt.Fprintf(w, `Registration successful! <a href="/sub">Go to Sub Page</a><script>setTimeout(() => window.location.href = '/sub', 5000);</script>`)
        }
    }
}


func adminHandler(w http.ResponseWriter, r *http.Request) {
    serveTemplate(w, "admin.html", nil)
}

func subHandler(w http.ResponseWriter, r *http.Request) {
    serveTemplate(w, "sub.html", nil)
}



func serveTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, data)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
    serveTemplate(w, "login.html", nil)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
    serveTemplate(w, "register.html", nil)
}

// Implement the new handlers here
func pagesHandler(w http.ResponseWriter, r *http.Request) {
	// List available pages
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	// Display a specific page
}

func requestAccessHandler(w http.ResponseWriter, r *http.Request) {
	// Handle access requests
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Display user status and tokens
}

func tradeStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Handle status token trading
}


func main() {
	var err error
	db, err = sql.Open("sqlite3", "./subad.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Routes
	r.Get("/", homeHandler)
	r.Get("/login", loginHandler)
	r.Post("/login", loginPostHandler)
	r.Get("/register", registerHandler)
	r.Post("/register", registerPostHandler)
	r.Get("/dashboard", dashboardHandler)
	r.Get("/admin", adminHandler)
	r.Get("/logout", logoutHandler)

	// New routes
	r.Get("/pages", pagesHandler)
	r.Get("/page/{id}", pageHandler)
	r.Post("/request-access", requestAccessHandler)
	r.Get("/status", statusHandler)
	r.Post("/trade-status", tradeStatusHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

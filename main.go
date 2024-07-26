package main
import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
	"strings"
	"fmt"

    _ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            is_admin BOOLEAN NOT NULL DEFAULT 0
        );
        CREATE TABLE IF NOT EXISTS pages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            content TEXT
        );
        CREATE TABLE IF NOT EXISTS page_access (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            page_id INTEGER,
            user_id INTEGER,
            FOREIGN KEY (page_id) REFERENCES pages(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)
    if err != nil {
        log.Fatal(err)
    }
}

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

func main() {
    db, err := sql.Open("sqlite3", "./subad.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    initDB(db)

    http.HandleFunc("/", loginPageHandler)
    http.HandleFunc("/register-page", registerPageHandler)
    http.HandleFunc("/login", loginHandler(db))
    http.HandleFunc("/register", registerHandler(db))
	http.HandleFunc("/admin", adminHandler)
    http.HandleFunc("/sub", subHandler)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

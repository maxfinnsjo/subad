package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"log"

	"github.com/maxfinnsjo/subad/database"
	"github.com/maxfinnsjo/subad/models"
	"github.com/maxfinnsjo/subad/sessions"
	"github.com/maxfinnsjo/subad/tokens"
)

type Handler struct {
	DB           *database.DB
	Sessions     *sessions.SessionStore
	TokenManager *tokens.TokenManager
}

func NewHandler(db *database.DB, sessions *sessions.SessionStore, tokenManager *tokens.TokenManager) *Handler {
    return &Handler{
        DB:           db,
        Sessions:     sessions,
        TokenManager: tokenManager,
    }
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Home",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/login.html"))
	err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Login",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.DB.GetUserByUsername(username)
	if err != nil || !user.CheckPassword(password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	session, err := h.Sessions.Create(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error creating session"})
		return
	}

	sessions.SetSession(w, session)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "redirect": "/dashboard"})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/register.html"))
	err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Register",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
    user := &models.User{
        Username: r.FormValue("username"),
        Email:    r.FormValue("email"),
        Password: r.FormValue("password"),
        Role:     "user",
    }
    err := h.DB.CreateUser(user)
    if err != nil {
        log.Printf("Error creating user: %v", err)
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }
    w.Write([]byte("User registered successfully"))
}


func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.DB.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/dashboard.html"))
	err = tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Dashboard",
		"User":  user,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err == nil {
		h.Sessions.Delete(sessionID.Value)
	}
	sessions.ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	page := &models.Page{
		Title:             r.FormValue("title"),
		Content:           r.FormValue("content"),
		OwnerID:           session.UserID,
		AccessLevel:       0,
		StatusRequirement: 0,
	}
	err = h.DB.CreatePage(page)
	if err != nil {
		http.Error(w, "Error creating page", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Page created successfully"))
}

func (h *Handler) ViewPage(w http.ResponseWriter, r *http.Request) {
	pageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}
	page, err := h.DB.GetPageByID(pageID)
	if err != nil {
		http.Error(w, "Error fetching page", http.StatusInternalServerError)
		return
	}

	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.DB.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}

	userStatus, err := h.TokenManager.CalculateUserStatus(user.ID)
	if err != nil {
		http.Error(w, "Error calculating user status", http.StatusInternalServerError)
		return
	}

	if !page.IsAccessibleBy(user, userStatus) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Write([]byte("Page Title: " + page.Title + "\nContent: " + page.Content))
}

func (h *Handler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	token, err := h.TokenManager.GenerateToken(session.UserID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Token generated successfully. Value: %d", token.Value)))
}

func (h *Handler) ViewUserStatus(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	status, err := h.TokenManager.CalculateUserStatus(session.UserID)
	if err != nil {
		http.Error(w, "Error calculating user status", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", status)))
}

func (h *Handler) EarnToken(w http.ResponseWriter, r *http.Request) {
	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	token, err := h.TokenManager.GenerateToken(session.UserID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Congratulations! You've earned a new token with value: %d", token.Value)))
}

func (h *Handler) TradeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := h.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenID, err := strconv.Atoi(r.FormValue("token_id"))
	if err != nil {
		http.Error(w, "Invalid token ID", http.StatusBadRequest)
		return
	}

	recipientID, err := strconv.Atoi(r.FormValue("recipient_id"))
	if err != nil {
		http.Error(w, "Invalid recipient ID", http.StatusBadRequest)
		return
	}

	err = h.TokenManager.TradeToken(session.UserID, recipientID, tokenID)
	if err != nil {
		http.Error(w, "Error trading token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Token traded successfully"))
}

func (h *Handler) getSession(r *http.Request) (*sessions.Session, error) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	session, ok := h.Sessions.Get(sessionID.Value)
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	return session, nil
}

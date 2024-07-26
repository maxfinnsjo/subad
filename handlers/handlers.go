package handlers

import (
	"net/http"
	"strconv"

	"your-project-path/database"
	"your-project-path/models"
)

type Handler struct {
	DB *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Implement home page logic
	w.Write([]byte("Welcome to Subad"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Implement login page logic
	w.Write([]byte("Login page"))
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	// Implement login post logic
	username := r.FormValue("username")
	password := r.FormValue("password")
	// TODO: Implement actual authentication logic
	w.Write([]byte("Login attempt for " + username))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Implement register page logic
	w.Write([]byte("Register page"))
}

func (h *Handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	// Implement register post logic
	user := &models.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Role:     "user",
	}
	err := h.DB.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Implement dashboard logic
	// TODO: Get actual user ID from session
	userID := 1
	user, err := h.DB.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Dashboard for " + user.Username))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic
	// TODO: Implement actual session destruction
	w.Write([]byte("Logged out successfully"))
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	// Implement page creation logic
	page := &models.Page{
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		OwnerID:     1, // TODO: Get actual user ID from session
		AccessLevel: 0, // TODO: Implement access level logic
	}
	err := h.DB.CreatePage(page)
	if err != nil {
		http.Error(w, "Error creating page", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Page created successfully"))
}

func (h *Handler) ViewPage(w http.ResponseWriter, r *http.Request) {
	// Implement page viewing logic
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
	// TODO: Check user permissions before displaying page
	w.Write([]byte("Page Title: " + page.Title + "\nContent: " + page.Content))
}

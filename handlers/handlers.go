package handlers

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/maxfinnsjo/subad/database"
	"github.com/maxfinnsjo/subad/models"
	"github.com/maxfinnsjo/subad/sessions"
	"github.com/maxfinnsjo/subad/tokens"
)

type Handler struct {
	DB           *database.DB
	Sessions     *sessions.SessionStore
	TokenManager *tokens.TokenManager
	Templates    *template.Template
}

func NewHandler(db *database.DB, sessions *sessions.SessionStore) *Handler {
	h := &Handler{
		DB:           db,
		Sessions:     sessions,
		TokenManager: tokens.NewTokenManager(db),
	}
	h.parseTemplates()
	return h
}

func (h *Handler) parseTemplates() {
	h.Templates = template.Must(template.ParseGlob("templates/*.html"))
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "home", map[string]interface{}{
		"Title": "Home",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "login", map[string]interface{}{
		"Title": "Login",
	})
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.DB.GetUserByUsername(username)
	if err != nil || !user.CheckPassword(password) {
		h.setFlashMessage(w, "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, err := h.Sessions.Create(user.ID)
	if err != nil {
		h.setFlashMessage(w, "Error creating session")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessions.SetSession(w, session)
	h.setFlashMessage(w, "Login successful")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "register", map[string]interface{}{
		"Title": "Register",
	})
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
		h.setFlashMessage(w, "Error creating user")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	h.setFlashMessage(w, "User registered successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := h.getUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	h.render(w, r, "dashboard", map[string]interface{}{
		"Title": "Dashboard",
		"User":  user,
	})
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	user := h.getUser(r)
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	h.render(w, r, "admin", map[string]interface{}{
		"Title": "Admin Dashboard",
		"User":  user,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err == nil {
		h.Sessions.Delete(sessionID.Value)
	}
	sessions.ClearSession(w)
	h.setFlashMessage(w, "Logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	data["FlashMessage"] = h.getFlashMessage(w, r)
	data["User"] = h.getUser(r)
	err := h.Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) setFlashMessage(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    url.QueryEscape(message),
		Path:     "/",
		HttpOnly: true,
	})
}

func (h *Handler) getFlashMessage(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("flash")
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		MaxAge: -1,
		Path:   "/",
	})
	message, _ := url.QueryUnescape(c.Value)
	return message
}

func (h *Handler) getUser(r *http.Request) *models.User {
	session, ok := h.getSession(r)
	if !ok {
		return nil
	}
	user, err := h.DB.GetUserByID(session.UserID)
	if err != nil {
		return nil
	}
	return user
}

func (h *Handler) getSession(r *http.Request) (*sessions.Session, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, false
	}
	return h.Sessions.Get(cookie.Value)
}

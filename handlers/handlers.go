package handlers

import (
	"net/http"
	"strconv"

	"your-project-path/database"
	"your-project-path/models"
	"your-project-path/sessions"
)

type Handler struct {
	DB      *database.DB
	Sessions *sessions.SessionStore

	func (h *Handler) GenerateToken(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, ok := h.Sessions.Get(sessionID.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenManager := tokens.NewTokenManager(h.DB)
		token, err := tokenManager.GenerateToken(session.UserID)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("Token generated successfully. Value: %d", token.Value)))
	}

	func (h *Handler) ViewUserStatus(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, ok := h.Sessions.Get(sessionID.Value)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenManager := tokens.NewTokenManager(h.DB)
		status, err := tokenManager.CalculateUserStatus(session.UserID)
		if err != nil {
			http.Error(w, "Error calculating user status", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("Your current status: %d", status)))
	}
}

func NewHandler(db *database.DB, sessions *sessions.SessionStore) *Handler {
	return &Handler{DB: db, Sessions: sessions}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Subad"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login page"))
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.DB.GetUserByUsername(username)
	if err != nil || !user.CheckPassword(password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := h.Sessions.Create(user.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	sessions.SetSession(w, session)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register page"))
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
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, ok := h.Sessions.Get(sessionID.Value)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.DB.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Dashboard for " + user.Username))
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
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, ok := h.Sessions.Get(sessionID.Value)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	page := &models.Page{
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		OwnerID:     session.UserID,
		AccessLevel: 0, // TODO: Implement access level logic
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
    
    sessionID, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    session, ok := h.Sessions.Get(sessionID.Value)
    if !ok {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    user, err := h.DB.GetUserByID(session.UserID)
    if err != nil {
        http.Error(w, "Error fetching user data", http.StatusInternalServerError)
        return
    }

    tokenManager := tokens.NewTokenManager(h.DB)
    userStatus, err := tokenManager.CalculateUserStatus(user.ID)
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

func (h *Handler) EarnToken(w http.ResponseWriter, r *http.Request) {
    sessionID, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    session, ok := h.Sessions.Get(sessionID.Value)
    if !ok {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    tokenManager := tokens.NewTokenManager(h.DB)
    token, err := tokenManager.GenerateToken(session.UserID)
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

    sessionID, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    session, ok := h.Sessions.Get(sessionID.Value)
    if !ok {
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

    tokenManager := tokens.NewTokenManager(h.DB)
    err = tokenManager.TradeToken(session.UserID, recipientID, tokenID)
    if err != nil {
        http.Error(w, "Error trading token: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("Token traded successfully"))
}

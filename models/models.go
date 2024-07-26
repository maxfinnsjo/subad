package models

import (
    "time"
)

type User struct {
    ID        int
    Username  string
    Email     string
    Password  string
    Role      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *User) CheckPassword(password string) bool {
    // For now, we'll do a simple comparison. In a real-world scenario,
    // you'd want to use a secure password hashing algorithm like bcrypt.
    return u.Password == password
}

type Page struct {
    ID               int
    Title            string
    Content          string
    OwnerID          int
    AccessLevel      int
    StatusRequirement int
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (p *Page) IsAccessibleBy(user *User, userStatus int) bool {
    if user.IsAdmin() || p.OwnerID == user.ID {
        return true
    }
    return userStatus >= p.StatusRequirement
}

type Subscription struct {
	ID        int
	UserID    int
	PageID    int
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StatusToken struct {
	ID        int
	UserID    int
	Value     int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (s *Subscription) IsActive() bool {
	return s.Status == "active" && s.ExpiresAt.After(time.Now())
}

func (st *StatusToken) IsValid() bool {
	return st.ExpiresAt.After(time.Now())
}

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

type Page struct {
	ID          int
	Title       string
	Content     string
	OwnerID     int
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

func (p *Page) IsAccessibleBy(user *User, subscriptions []Subscription) bool {
	if user.IsAdmin() || p.OwnerID == user.ID {
		return true
	}

	for _, sub := range subscriptions {
		if sub.PageID == p.ID && sub.UserID == user.ID && sub.Status == "active" && sub.ExpiresAt.After(time.Now()) {
			return true
		}
	}

	return false
}

func (s *Subscription) IsActive() bool {
	return s.Status == "active" && s.ExpiresAt.After(time.Now())
}

func (st *StatusToken) IsValid() bool {
	return st.ExpiresAt.After(time.Now())
}

package database

import (
    "database/sql"
    "fmt"
    "time"
    "github.com/maxfinnsjo/subad/models"
)

type DB struct {
    *sql.DB
}

func (db *DB) MarkTokenAsUsed(tokenID int) error {
    // Implement the logic to mark a token as used
    return nil
}

func (db *DB) GetTokenByID(tokenID int) (*models.StatusToken, error) {
    // Implement the logic to retrieve a token by ID
    return nil, nil
}

func (db *DB) UpdateTokenOwner(tokenID, newOwnerID int) error {
    // Implement the logic to update a token's owner
    return nil

	func (db *DB) CreateToken(token *models.StatusToken) error {
		_, err := db.Exec("INSERT INTO status_tokens (user_id, value, created_at, expires_at) VALUES (?, ?, ?, ?)",
			token.UserID, token.Value, token.CreatedAt, token.ExpiresAt)
		if err != nil {
			return fmt.Errorf("error creating token: %w", err)
		}
		return nil
	}
	

}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &DB{DB: db}, nil
}

func (db *DB) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow("SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (db *DB) CreateUser(user *models.User) error {
	_, err := db.Exec("INSERT INTO users (username, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.Role, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (db *DB) GetPageByID(id int) (*models.Page, error) {
    page := &models.Page{}
    err := db.QueryRow("SELECT id, title, content, owner_id, access_level, status_requirement, created_at, updated_at FROM pages WHERE id = ?", id).
        Scan(&page.ID, &page.Title, &page.Content, &page.OwnerID, &page.AccessLevel, &page.StatusRequirement, &page.CreatedAt, &page.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("error getting page: %w", err)
    }
    return page, nil
}

func (db *DB) CreatePage(page *models.Page) error {
    _, err := db.Exec("INSERT INTO pages (title, content, owner_id, access_level, status_requirement, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
        page.Title, page.Content, page.OwnerID, page.AccessLevel, page.StatusRequirement, time.Now(), time.Now())
    if err != nil {
        return fmt.Errorf("error creating page: %w", err)
    }
    return nil
}

func (db *DB) GetSubscriptionsByUserID(userID int) ([]models.Subscription, error) {
	rows, err := db.Query("SELECT id, user_id, page_id, status, expires_at, created_at, updated_at FROM subscriptions WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.PageID, &sub.Status, &sub.ExpiresAt, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (db *DB) CreateSubscription(sub *models.Subscription) error {
	_, err := db.Exec("INSERT INTO subscriptions (user_id, page_id, status, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		sub.UserID, sub.PageID, sub.Status, sub.ExpiresAt, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	return nil
}

func (db *DB) GetStatusTokensByUserID(userID int) ([]models.StatusToken, error) {
	rows, err := db.Query("SELECT id, user_id, value, created_at, expires_at FROM status_tokens WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting status tokens: %w", err)
	}
	defer rows.Close()

	var tokens []models.StatusToken
	for rows.Next() {
		var token models.StatusToken
		err := rows.Scan(&token.ID, &token.UserID, &token.Value, &token.CreatedAt, &token.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning status token: %w", err)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (db *DB) CreateStatusToken(token *models.StatusToken) error {
	_, err := db.Exec("INSERT INTO status_tokens (user_id, value, created_at, expires_at) VALUES (?, ?, ?, ?)",
		token.UserID, token.Value, time.Now(), token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error creating status token: %w", err)
	}
	return nil
}


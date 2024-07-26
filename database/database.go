package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	*sql.DB
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

// Add methods for database operations here, for example:

func (db *DB) GetUserByID(id int) (/* User struct */, error) {
	// Implement user retrieval logic
}

func (db *DB) CreateUser(/* User struct */) error {
	// Implement user creation logic
}

// Add more methods as needed

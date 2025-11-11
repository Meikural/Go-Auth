package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// InitDB initializes and returns a database connection
func InitDB(driver, source string) (*sql.DB, error) {
	database, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(5)

	return database, nil
}
package database

import (
	"database/sql"
	"fmt"
)

// CreateTables creates the necessary tables in the database
func CreateTables(db *sql.DB) error {
	// Create roles table
	createRolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createRolesTable)
	if err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	// Create users table with role reference
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(100) DEFAULT 'User',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
	"time"
)

// CreateUser inserts a new user into the database with role
func CreateUser(db *sql.DB, username, email, hashedPassword, role string) (*models.User, error) {
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
	INSERT INTO users (username, email, password, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	err := db.QueryRow(query, username, email, hashedPassword, role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
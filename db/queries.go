package db

import (
	"database/sql"
	"errors"
	"fmt"
	"go-auth/models"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidPassword  = errors.New("invalid password")
)

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, username, email, hashedPassword string) (*models.User, error) {
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := db.QueryRow(query, username, email, hashedPassword, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email address
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdatePassword updates a user's password
func UpdatePassword(db *sql.DB, userID int, hashedPassword string) error {
	query := `
	UPDATE users
	SET password = $1, updated_at = $2
	WHERE id = $3
	`

	result, err := db.Exec(query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
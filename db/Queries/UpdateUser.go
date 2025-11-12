package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
	"time"
)

// UpdateUser updates a user's username and/or email
func UpdateUser(db *sql.DB, userID string, username, email string) (*models.User, error) {
	user := &models.User{}

	query := `
	UPDATE users
	SET username = $1, email = $2, updated_at = $3
	WHERE id = $4 AND deleted_at IS NULL
	RETURNING id, username, email, password, role, deleted_at, created_at, updated_at
	`

	err := db.QueryRow(query, username, email, time.Now(), userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.DeletedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		// Check if it's a unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	user.Password = "" // Don't expose password
	return user, nil
}
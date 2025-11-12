package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
	"time"
)

// UpdateUserRole updates a user's role
func UpdateUserRole(db *sql.DB, userID string, newRole string) (*models.User, error) {
	user := &models.User{}

	query := `
	UPDATE users
	SET role = $1, updated_at = $2
	WHERE id = $3
	RETURNING id, username, email, password, role, created_at, updated_at
	`

	err := db.QueryRow(query, newRole, time.Now(), userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	user.Password = "" // Don't expose password
	return user, nil
}
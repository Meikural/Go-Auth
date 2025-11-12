package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// GetUserByEmail retrieves a user by email address (excludes soft-deleted users)
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, role, deleted_at, created_at, updated_at
	FROM users
	WHERE email = $1 AND deleted_at IS NULL
	`

	err := db.QueryRow(query, email).Scan(
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
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
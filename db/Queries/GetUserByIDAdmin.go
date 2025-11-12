package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// GetUserByIDAdmin retrieves a user by ID for admin (includes deleted users info)
func GetUserByIDAdmin(db *sql.DB, id string) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, role, deleted_at, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(
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

	user.Password = "" // Don't expose password
	return user, nil
}
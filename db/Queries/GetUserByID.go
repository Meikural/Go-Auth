package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(
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
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

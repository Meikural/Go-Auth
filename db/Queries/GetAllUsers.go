package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *sql.DB) ([]*models.User, error) {
	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Password = "" // Don't expose passwords
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
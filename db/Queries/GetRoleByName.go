package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// GetRoleByName retrieves a role by name
func GetRoleByName(db *sql.DB, roleName string) (*models.Role, error) {
	role := &models.Role{}

	query := `
	SELECT id, name
	FROM roles
	WHERE name = $1
	`

	err := db.QueryRow(query, roleName).Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}
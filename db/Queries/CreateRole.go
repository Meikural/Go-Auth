package queries

import (
	"database/sql"
	"fmt"
	"go-auth/models"
)

// CreateRole inserts a new role into the database
func CreateRole(db *sql.DB, roleName string) (*models.Role, error) {
	role := &models.Role{
		Name: roleName,
	}

	query := `
	INSERT INTO roles (name)
	VALUES ($1)
	ON CONFLICT (name) DO NOTHING
	RETURNING id
	`

	err := db.QueryRow(query, roleName).Scan(&role.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Role already exists, fetch it
			return GetRoleByName(db, roleName)
		}
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}
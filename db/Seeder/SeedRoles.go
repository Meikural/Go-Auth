package seeder

import (
	"database/sql"
	"fmt"
	queries "go-auth/db/Queries"
)

// SeedRoles creates all roles from the provided list
func SeedRoles(db *sql.DB, roles []string) error {
	for _, roleName := range roles {
		_, err := queries.CreateRole(db, roleName)
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %w", roleName, err)
		}
	}
	return nil
}
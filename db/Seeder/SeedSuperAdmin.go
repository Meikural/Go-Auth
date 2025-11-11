package seeder

import (
	"database/sql"
	"fmt"
	queries "go-auth/db/Queries"
	passw "go-auth/utils/password"
)

// SeedSuperAdmin creates a super admin user if it doesn't exist
func SeedSuperAdmin(db *sql.DB, email, password, superAdminRole string) error {
	// Check if super admin already exists
	existingUser, err := queries.GetUserByEmail(db, email)
	if err == nil && existingUser != nil {
		// User already exists, skip creation
		fmt.Printf("Super admin user already exists: %s\n", email)
		return nil
	}

	if err != nil && err != queries.ErrUserNotFound {
		return fmt.Errorf("failed to check for existing super admin: %w", err)
	}

	// Hash the password
	hashedPassword, err := passw.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash super admin password: %w", err)
	}

	// Create super admin user
	_, err = queries.CreateUser(db, "superadmin", email, hashedPassword, superAdminRole)
	if err != nil {
		if err == queries.ErrUserExists {
			fmt.Printf("Super admin user already exists\n")
			return nil
		}
		return fmt.Errorf("failed to create super admin user: %w", err)
	}

	fmt.Printf("Super admin user created successfully: %s\n", email)
	return nil
}
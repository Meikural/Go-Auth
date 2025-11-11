package db

import (
	"database/sql"
	"fmt"
	"go-auth/utils"
)

// SeedRoles creates all roles from the provided list
func SeedRoles(db *sql.DB, roles []string) error {
	for _, roleName := range roles {
		_, err := CreateRole(db, roleName)
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %w", roleName, err)
		}
	}
	return nil
}

// SeedSuperAdmin creates a super admin user if it doesn't exist
func SeedSuperAdmin(db *sql.DB, email, password, superAdminRole string) error {
	// Check if super admin already exists
	existingUser, err := GetUserByEmail(db, email)
	if err == nil && existingUser != nil {
		// User already exists, skip creation
		fmt.Printf("Super admin user already exists: %s\n", email)
		return nil
	}

	if err != nil && err != ErrUserNotFound {
		return fmt.Errorf("failed to check for existing super admin: %w", err)
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash super admin password: %w", err)
	}

	// Create super admin user
	_, err = CreateUser(db, "superadmin", email, hashedPassword, superAdminRole)
	if err != nil {
		if err == ErrUserExists {
			fmt.Printf("Super admin user already exists\n")
			return nil
		}
		return fmt.Errorf("failed to create super admin user: %w", err)
	}

	fmt.Printf("Super admin user created successfully: %s\n", email)
	return nil
}
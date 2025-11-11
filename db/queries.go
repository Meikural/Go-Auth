package db

import (
	"database/sql"
	"errors"
	"fmt"
	"go-auth/models"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrRoleNotFound     = errors.New("role not found")
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

// CreateUser inserts a new user into the database with role
func CreateUser(db *sql.DB, username, email, hashedPassword, role string) (*models.User, error) {
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
	INSERT INTO users (username, email, password, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	err := db.QueryRow(query, username, email, hashedPassword, role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email address
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	query := `
	SELECT id, username, email, password, role, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	err := db.QueryRow(query, email).Scan(
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

// UpdatePassword updates a user's password
func UpdatePassword(db *sql.DB, userID int, hashedPassword string) error {
	query := `
	UPDATE users
	SET password = $1, updated_at = $2
	WHERE id = $3
	`

	result, err := db.Exec(query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

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

// UpdateUserRole updates a user's role
func UpdateUserRole(db *sql.DB, userID int, newRole string) (*models.User, error) {
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
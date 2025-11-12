package queries

import (
	"database/sql"
	"fmt"
	"time"
)

// UpdatePassword updates a user's password
func UpdatePassword(db *sql.DB, userID string, hashedPassword string) error {
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
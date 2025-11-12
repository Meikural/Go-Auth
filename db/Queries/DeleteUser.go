package queries

import (
	"database/sql"
	"fmt"
	"time"
)

// DeleteUser performs a soft delete on a user by setting deleted_at timestamp
func DeleteUser(db *sql.DB, userID string) error {
	query := `
	UPDATE users
	SET deleted_at = $1, updated_at = $2
	WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := db.Exec(query, time.Now(), time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
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
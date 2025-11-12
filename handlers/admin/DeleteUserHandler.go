package admin

import (
	"database/sql"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/middleware/auth"
	"net/http"
	"strings"
)

// DeleteUserHandler soft deletes a user (Super Admin only)
func DeleteUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get claims from context (set by middleware)
		claims, err := auth.GetClaimsFromContext(r)
		if err != nil {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		// Extract user ID from URL path
		// Expected format: /admin/users/delete/{uuid}
		userID := strings.TrimPrefix(r.URL.Path, "/admin/users/delete/")
		if userID == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		// Prevent self-deletion
		if userID == claims.UserID {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot delete your own account"})
			return
		}

		// Soft delete the user
		err = queries.DeleteUser(database, userID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete user"})
			return
		}

		handlers.RespondJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
	}
}
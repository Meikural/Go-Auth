package admin

import (
	"database/sql"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"net/http"
	"strings"
)

// GetUserHandler returns a specific user's details (Super Admin only)
func GetUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Extract user ID from URL path
		// Expected format: /admin/users/get/{uuid}
		userID := strings.TrimPrefix(r.URL.Path, "/admin/users/get/")
		if userID == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		// Get user from database (admin can see deleted users)
		user, err := queries.GetUserByIDAdmin(database, userID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		handlers.RespondJSON(w, http.StatusOK, user)
	}
}
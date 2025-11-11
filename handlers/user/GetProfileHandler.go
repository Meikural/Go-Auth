package user

import (
	"database/sql"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/middleware/auth"
	"net/http"
)

// GetProfileHandler returns the current authenticated user's profile
func GetProfileHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get claims from context (set by middleware)
		claims, err := auth.GetClaimsFromContext(r)
		if err != nil {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		// Get user from database
		user, err := queries.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Don't expose password in response
		user.Password = ""

		handlers.RespondJSON(w, http.StatusOK, user)
	}
}
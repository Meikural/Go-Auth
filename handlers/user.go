package handlers

import (
	"database/sql"
	"go-auth/db"
	"go-auth/middleware"
	"net/http"
)

// GetProfileHandler returns the current authenticated user's profile
func GetProfileHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get claims from context (set by middleware)
		claims, err := middleware.GetClaimsFromContext(r)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		// Get user from database
		user, err := db.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == db.ErrUserNotFound {
				respondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Don't expose password in response
		user.Password = ""

		respondJSON(w, http.StatusOK, user)
	}
}

// HealthCheckHandler is a simple endpoint to check if the service is running
func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	}
}
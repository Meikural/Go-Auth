package auth

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/middleware/auth"
	"go-auth/models"
	"go-auth/utils/password"
	"net/http"
)

// ChangePasswordHandler handles password changes (requires authentication)
func ChangePasswordHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get claims from context (set by middleware)
		claims, err := auth.GetClaimsFromContext(r)
		if err != nil {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		var req models.ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.OldPassword == "" || req.NewPassword == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "old password and new password are required"})
			return
		}

		// Get user
		user, err := queries.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Verify old password
		if !password.VerifyPassword(user.Password, req.OldPassword) {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
			return
		}

		// Hash new password
		hashedPassword, err := password.HashPassword(req.NewPassword)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		// Update password
		if err := queries.UpdatePassword(database, claims.UserID, hashedPassword); err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update password"})
			return
		}

		handlers.RespondJSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
	}
}
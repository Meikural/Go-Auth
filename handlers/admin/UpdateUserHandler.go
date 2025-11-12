package admin

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/models"
	"net/http"
	"strings"
)

// UpdateUserHandler updates a user's username and/or email (Super Admin only)
func UpdateUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Extract user ID from URL path
		// Expected format: /admin/users/update/{uuid}
		userID := strings.TrimPrefix(r.URL.Path, "/admin/users/update/")
		if userID == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		// Parse request body
		var req models.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate at least one field is provided
		if req.Username == "" && req.Email == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "at least username or email must be provided"})
			return
		}

		// Get current user to preserve fields not being updated
		currentUser, err := queries.GetUserByIDAdmin(database, userID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Use provided values or keep existing ones
		username := req.Username
		if username == "" {
			username = currentUser.Username
		}

		email := req.Email
		if email == "" {
			email = currentUser.Email
		}

		// Update user
		updatedUser, err := queries.UpdateUser(database, userID, username, email)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			if err == queries.ErrUserExists {
				handlers.RespondJSON(w, http.StatusConflict, map[string]string{"error": "username or email already exists"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update user"})
			return
		}

		handlers.RespondJSON(w, http.StatusOK, updatedUser)
	}
}
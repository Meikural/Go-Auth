package admin

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/middleware/auth"
	"go-auth/models"
	"net/http"
	"strings"
)

// UpdateUserRoleHandler updates a user's role (Super Admin only)
func UpdateUserRoleHandler(database *sql.DB, availableRoles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
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
		// Expected format: /admin/users/uuid-string/role
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid path"})
			return
		}

		userID := pathParts[3]
		if userID == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		// Parse request body
		var req models.UpdateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate role is not empty
		if req.Role == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "role is required"})
			return
		}

		// Validate role exists in available roles
		roleExists := false
		for _, role := range availableRoles {
			if strings.EqualFold(role, req.Role) {
				req.Role = role // Use the correct casing
				roleExists = true
				break
			}
		}
		if !roleExists {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role"})
			return
		}

		// Prevent updating own role (optional but safer)
		if userID == claims.UserID {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot change your own role"})
			return
		}

		// Get current user to verify they're not removing the last super admin
		targetUser, err := queries.GetUserByID(database, userID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Prevent removing super admin role if they're the last one
		if strings.EqualFold(targetUser.Role, availableRoles[0]) && !strings.EqualFold(req.Role, availableRoles[0]) {
			// Check if there are other super admins
			allUsers, _ := queries.GetAllUsers(database)
			superAdminCount := 0
			for _, u := range allUsers {
				if strings.EqualFold(u.Role, availableRoles[0]) {
					superAdminCount++
				}
			}
			if superAdminCount <= 1 {
				handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot remove the last super admin"})
				return
			}
		}

		// Update user role
		updatedUser, err := queries.UpdateUserRole(database, userID, req.Role)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update user role"})
			return
		}

		response := models.UpdateRoleResponse{
			Message: "user role updated successfully",
			User:    updatedUser,
		}

		handlers.RespondJSON(w, http.StatusOK, response)
	}
}
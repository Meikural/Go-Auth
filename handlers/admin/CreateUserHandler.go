package admin

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/utils/password"
	"net/http"
)

// CreateUserRequest is the payload for creating a new user by admin
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// CreateUserHandler creates a new user (Super Admin only)
func CreateUserHandler(database *sql.DB, availableRoles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Parse request body
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email, password, and role are required"})
			return
		}

		// Validate role exists in available roles
		roleExists := false
		for _, role := range availableRoles {
			if role == req.Role {
				roleExists = true
				break
			}
		}
		if !roleExists {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role"})
			return
		}

		// Hash password
		hashedPassword, err := password.HashPassword(req.Password)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		// Create user
		user, err := queries.CreateUser(database, req.Username, req.Email, hashedPassword, req.Role)
		if err != nil {
			if err == queries.ErrUserExists {
				handlers.RespondJSON(w, http.StatusConflict, map[string]string{"error": "user already exists"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
			return
		}

		handlers.RespondJSON(w, http.StatusCreated, user)
	}
}
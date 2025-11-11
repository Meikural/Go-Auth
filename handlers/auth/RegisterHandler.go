package auth

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/models"
	"go-auth/utils/jwt"
	"go-auth/utils/password"
	"net/http"
)

// RegisterHandler handles user registration
func RegisterHandler(database *sql.DB, secretKey, defaultRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.Username == "" || req.Email == "" || req.Password == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email, and password are required"})
			return
		}

		// Hash password
		hashedPassword, err := password.HashPassword(req.Password)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		// Create user with default role
		user, err := queries.CreateUser(database, req.Username, req.Email, hashedPassword, defaultRole)
		if err != nil {
			if err == queries.ErrUserExists {
				handlers.RespondJSON(w, http.StatusConflict, map[string]string{"error": "user already exists"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
			return
		}

		// Generate tokens
		accessToken, err := jwt.GenerateToken(user.ID, user.Username, user.Email, user.Role, models.AccessToken, secretKey)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		refreshToken, err := jwt.GenerateToken(user.ID, user.Username, user.Email, user.Role, models.RefreshToken, secretKey)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate refresh token"})
			return
		}

		// Don't expose password in response
		user.Password = ""

		response := models.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         *user,
		}

		handlers.RespondJSON(w, http.StatusCreated, response)
	}
}
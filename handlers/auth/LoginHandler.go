package auth

import (
	"database/sql"
	"encoding/json"
	"go-auth/db"
	"go-auth/handlers"
	"go-auth/models"
	"go-auth/utils"
	"net/http"
)

// LoginHandler handles user login
func LoginHandler(database *sql.DB, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
			return
		}

		// Get user by email
		user, err := db.GetUserByEmail(database, req.Email)
		if err != nil {
			if err == db.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Verify password
		if !utils.VerifyPassword(user.Password, req.Password) {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		// Generate tokens
		accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role, models.AccessToken, secretKey)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		refreshToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role, models.RefreshToken, secretKey)
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

		handlers.RespondJSON(w, http.StatusOK, response)
	}
}
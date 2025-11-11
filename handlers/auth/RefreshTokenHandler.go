package auth

import (
	"database/sql"
	"encoding/json"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/models"
	"go-auth/utils/jwt"
	"net/http"
)

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(database *sql.DB, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.RefreshToken == "" {
			handlers.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh token is required"})
			return
		}

		// Verify refresh token
		claims, err := jwt.VerifyToken(req.RefreshToken, secretKey)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token expired"})
				return
			}
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
			return
		}

		// Verify it's a refresh token
		if claims.TokenType != models.RefreshToken {
			handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token type"})
			return
		}

		// Get user to ensure they still exist
		user, err := queries.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == queries.ErrUserNotFound {
				handlers.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
				return
			}
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Generate new access token
		accessToken, err := jwt.GenerateToken(user.ID, user.Username, user.Email, user.Role, models.AccessToken, secretKey)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		response := map[string]string{
			"access_token": accessToken,
		}

		handlers.RespondJSON(w, http.StatusOK, response)
	}
}
package handlers

import (
	"database/sql"
	"encoding/json"
	"go-auth/db"
	"go-auth/middleware"
	"go-auth/models"
	"go-auth/utils"
	"net/http"
)

// RegisterHandler handles user registration
func RegisterHandler(database *sql.DB, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.Username == "" || req.Email == "" || req.Password == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email, and password are required"})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		// Create user
		user, err := db.CreateUser(database, req.Username, req.Email, hashedPassword)
		if err != nil {
			if err == db.ErrUserExists {
				respondJSON(w, http.StatusConflict, map[string]string{"error": "user already exists"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
			return
		}

		// Generate tokens
		accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, models.AccessToken, secretKey)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		refreshToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, models.RefreshToken, secretKey)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate refresh token"})
			return
		}

		// Don't expose password in response
		user.Password = ""

		response := models.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         *user,
		}

		respondJSON(w, http.StatusCreated, response)
	}
}

// LoginHandler handles user login
func LoginHandler(database *sql.DB, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.Email == "" || req.Password == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
			return
		}

		// Get user by email
		user, err := db.GetUserByEmail(database, req.Email)
		if err != nil {
			if err == db.ErrUserNotFound {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Verify password
		if !utils.VerifyPassword(user.Password, req.Password) {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		// Generate tokens
		accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, models.AccessToken, secretKey)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		refreshToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, models.RefreshToken, secretKey)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate refresh token"})
			return
		}

		// Don't expose password in response
		user.Password = ""

		response := models.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         *user,
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(database *sql.DB, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req models.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if req.RefreshToken == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh token is required"})
			return
		}

		// Verify refresh token
		claims, err := utils.VerifyToken(req.RefreshToken, secretKey)
		if err != nil {
			if err == utils.ErrExpiredToken {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token expired"})
				return
			}
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
			return
		}

		// Verify it's a refresh token
		if claims.TokenType != models.RefreshToken {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token type"})
			return
		}

		// Get user to ensure they still exist
		user, err := db.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == db.ErrUserNotFound {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Generate new access token
		accessToken, err := utils.GenerateToken(user.ID, user.Username, user.Email, models.AccessToken, secretKey)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
			return
		}

		response := map[string]string{
			"access_token": accessToken,
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// ChangePasswordHandler handles password changes (requires authentication)
func ChangePasswordHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get claims from context (set by middleware)
		claims, err := middleware.GetClaimsFromContext(r)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		var req models.ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		// Validate input
		if req.OldPassword == "" || req.NewPassword == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "old password and new password are required"})
			return
		}

		// Get user
		user, err := db.GetUserByID(database, claims.UserID)
		if err != nil {
			if err == db.ErrUserNotFound {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
			return
		}

		// Verify old password
		if !utils.VerifyPassword(user.Password, req.OldPassword) {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
			return
		}

		// Hash new password
		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}

		// Update password
		if err := db.UpdatePassword(database, claims.UserID, hashedPassword); err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update password"})
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
	}
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
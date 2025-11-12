package auth

import (
	"context"
	"go-auth/middleware/constants"
	"go-auth/models"
	"go-auth/utils/jwt"
	"net/http"
	"strings"
)

// AuthMiddleware verifies JWT token and extracts claims
func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				constants.RespondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Bearer token format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				constants.RespondError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Verify token
			claims, err := jwt.VerifyToken(tokenString, secretKey)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					constants.RespondError(w, http.StatusUnauthorized, "token expired")
					return
				}
				constants.RespondError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Verify it's an access token
			if claims.TokenType != models.AccessToken {
				constants.RespondError(w, http.StatusUnauthorized, "invalid token type")
				return
			}

			// Store claims in context for handler access
			ctx := context.WithValue(r.Context(), constants.UserContextKey, claims)
			r = r.WithContext(ctx)

			// Also store in headers for easier access
			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-Username", claims.Username)
			r.Header.Set("X-Email", claims.Email)
			r.Header.Set("X-Role", claims.Role)

			next.ServeHTTP(w, r)
		})
	}
}
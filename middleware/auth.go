package middleware

import (
	"context"
	"fmt"
	"go-auth/models"
	"go-auth/utils"
	"net/http"
	"strings"
)

// ContextKey is used to store values in request context
type ContextKey string

const UserContextKey ContextKey = "user"

// AuthMiddleware verifies JWT token and extracts claims
func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			// Bearer token format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Verify token
			claims, err := utils.VerifyToken(tokenString, secretKey)
			if err != nil {
				if err == utils.ErrExpiredToken {
					respondError(w, http.StatusUnauthorized, "token expired")
					return
				}
				respondError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Verify it's an access token
			if claims.TokenType != models.AccessToken {
				respondError(w, http.StatusUnauthorized, "invalid token type")
				return
			}

			// Store claims in context for handler access
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(r *http.Request) (*models.Claims, error) {
	claims, ok := r.Context().Value(UserContextKey).(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

// respondError writes a JSON error response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"error":"%s"}`, message)
}
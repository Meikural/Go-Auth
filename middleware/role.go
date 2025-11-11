package middleware

import (
	"fmt"
	"go-auth/middleware/constants"

	"net/http"
)

// RequireRole middleware checks if user has the required role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from headers (set by AuthMiddleware)
			userRole := r.Header.Get("X-Role")
			if userRole == "" {
				constants.RespondError(w, http.StatusUnauthorized, "user role not found")
				return
			}

			// Check if user has required role
			if userRole != requiredRole {
				constants.RespondError(w, http.StatusForbidden, fmt.Sprintf("this action requires %s role", requiredRole))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
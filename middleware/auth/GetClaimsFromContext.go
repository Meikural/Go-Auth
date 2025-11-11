package auth

import (
	"fmt"
	"go-auth/middleware/constants"
	"go-auth/models"
	"net/http"
)

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(r *http.Request) (*models.Claims, error) {
	claims, ok := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
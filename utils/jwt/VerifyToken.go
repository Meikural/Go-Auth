package jwt

import (
	"fmt"
	"go-auth/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken validates a JWT token and returns the claims
func VerifyToken(tokenString, secretKey string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	// Check if token is valid
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token has expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
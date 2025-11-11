package jwt

import (
	"go-auth/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token with the given claims
func GenerateToken(userID int, username, email, role string, tokenType models.TokenType, secretKey string) (string, error) {
	var expirationTime time.Time

	// Set expiration based on token type
	if tokenType == models.AccessToken {
		expirationTime = time.Now().Add(models.AccessTokenDuration)
	} else {
		expirationTime = time.Now().Add(models.RefreshTokenDuration)
	}

	claims := &models.Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
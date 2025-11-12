package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType represents the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenDuration constants
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days
)

// Claims represents the JWT claims
type Claims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Token represents a token response
type Token struct {
	Value     string    `json:"value"`
	Type      TokenType `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}
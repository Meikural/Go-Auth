package jwt

import "errors"

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token expired")
	ErrInvalidClaims   = errors.New("invalid claims")
)
package token

import (
	"errors"

	"github.com/go-jose/go-jose/v4/jwt"
)

var (
	ErrInvalidKey = errors.New("invalid key length, must be 32 bytes")
)

// DataWithClaims holds user data with jwt.Claims
type DataWithClaims interface {
	GetClaims() *jwt.Claims
}

// EmbedJWTClaims contains jwt.Claims
// Embed this in user structures to use with Encoder and Decoder interfaces
type EmbedJWTClaims struct {
	jwt.Claims
}

// GetClaims returns the embedded jwt.Claims
func (e *EmbedJWTClaims) GetClaims() *jwt.Claims {
	return &e.Claims
}

package ports

import (
	"time"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
)

type TokenClaims map[string]any

const (
	// AccessTokenLifetime  = 1 * time.Minute
	AccessTokenLifetime  = 15 * time.Minute
	RefreshTokenLifetime = 1 * 24 * time.Hour
)

type TokenManager interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken(user *domain.User) (string, error)
	RefreshAccessToken(user *domain.User, refreshToken string) (string, error)
	ValidateToken(tokenString string) (TokenClaims, error)
}

package ports

import "github.com/guilehm/echo-vision/echo-hub/internal/app/domain"

type TokenClaims map[string]any

type TokenManager interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken(user *domain.User) (string, error)
	RefreshAccessToken(user *domain.User, refreshToken string) (string, error)
	ValidateToken(tokenString string) (TokenClaims, error)
}

package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type UserPort interface {
	CreateUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) (uuid.UUID, error)
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
	UserByAccessToken(ctx context.Context, accessToken string) (*domain.User, error)
	UserByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
}

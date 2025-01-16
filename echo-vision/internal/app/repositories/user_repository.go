package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type UserRepository interface {
	SaveUser(ctx context.Context, tx Transaction, user *domain.User) (uuid.UUID, error)
	UpdateTokens(ctx context.Context, tx Transaction, accessToken, refreshToken string, userID uuid.UUID) error
	FindUserByEmail(ctx context.Context, tx Transaction, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, tx Transaction, id uuid.UUID) (*domain.User, error)
	FindUserByTokens(ctx context.Context, tx Transaction, accessToken, refreshToken string) (*domain.User, error)
}

package repositories

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/domain"
)

type UserRepository interface {
	SaveUser(ctx context.Context, tx Transaction, event *domain.User) error
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

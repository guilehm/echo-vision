package repositories

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/domain"
)

type UserRepository interface {
	Save(ctx context.Context, tx Transaction, event *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

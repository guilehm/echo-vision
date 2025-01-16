package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type UserPort interface {
	FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	SaveUser(ctx context.Context, user *domain.User) (uuid.UUID, error)
}

type UserCreateInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

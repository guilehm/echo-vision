package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/repositories"
)

type ManageUsers struct {
	Repository repositories.Repository
}

func (uc *ManageUsers) FindUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	return uc.Repository.FindUserByID(ctx, nil, id)
}

func (uc *ManageUsers) FindUserByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	return uc.Repository.FindUserByEmail(ctx, nil, email)
}

func (uc *ManageUsers) SaveUser(
	ctx context.Context,
	user *domain.User,
) (uuid.UUID, error) {
	return uc.Repository.SaveUser(ctx, nil, user)
}

func NewManageUsersUseCase(repository repositories.Repository) ports.UserPort {
	return &ManageUsers{
		Repository: repository,
	}
}

package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/repositories"
)

type ManageUsers struct {
	repository      repositories.Repository
	tokenManager    ports.TokenManager
	passwordManager ports.PasswordManager
}

func (uc *ManageUsers) FindUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	return uc.repository.FindUserByID(ctx, nil, id)
}

func (uc *ManageUsers) FindUserByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	return uc.repository.FindUserByEmail(ctx, nil, email)
}

func (uc *ManageUsers) CreateUser(
	ctx context.Context,
	firstName, lastName, email, password string,
) (*domain.User, error) {
	now := time.Now()
	user := domain.NewUser(
		uuid.New(),
		firstName,
		lastName,
		email,
		now,
		now,
	)
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := user.IsValidPassword(password); err != nil {
		return nil, err
	}
	fmt.Println("us.passwordManager", uc.passwordManager)

	hp, err := uc.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.SetHashedPassword(hp)

	accessToken, err := uc.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.tokenManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	user.SetTokens(accessToken, refreshToken)

	return user, nil
}

func (uc *ManageUsers) SaveUser(
	ctx context.Context,
	user *domain.User,
) (uuid.UUID, error) {
	if err := user.Validate(); err != nil {
		return uuid.Nil, err
	}
	return uc.repository.SaveUser(ctx, nil, user)
}

func NewManageUsersUseCase(
	repository repositories.Repository,
	tokenManager ports.TokenManager,
	passwordManager ports.PasswordManager,
) ports.UserPort {
	return &ManageUsers{
		repository:      repository,
		tokenManager:    tokenManager,
		passwordManager: passwordManager,
	}
}

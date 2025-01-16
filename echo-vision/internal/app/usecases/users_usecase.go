package usecases

import (
	"context"
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

// AuthenticateUser implements ports.UserPort.
func (uc *ManageUsers) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	var u *domain.User
	err := uc.repository.WithTransaction(ctx, func(ctx context.Context, tx repositories.Transaction) error {
		user, err := uc.repository.FindUserByEmail(ctx, tx, email)
		if err != nil {
			return err
		}
		if err := uc.passwordManager.ValidatePassword(password, user.Password()); err != nil {
			return err
		}
		accessToken, err := uc.tokenManager.GenerateAccessToken(user)
		if err != nil {
			return err
		}

		refreshToken, err := uc.tokenManager.GenerateRefreshToken(user)
		if err != nil {
			return err
		}

		user.SetTokens(accessToken, refreshToken)
		err = uc.repository.UpdateTokens(ctx, tx, accessToken, refreshToken, user.ID())
		if err != nil {
			return err
		}
		u = user
		return nil
	})
	return u, err
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

// UserByAccessToken implements ports.UserPort.
func (uc *ManageUsers) UserByAccessToken(ctx context.Context, accessToken string) (*domain.User, error) {
	u, err := uc.repository.FindUserByTokens(ctx, nil, accessToken, "")
	if err != nil {
		return nil, err
	}
	_, err = uc.tokenManager.ValidateToken(u.AccessToken())
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UserByRefreshToken implements ports.UserPort.
func (uc *ManageUsers) UserByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	u, err := uc.repository.FindUserByTokens(ctx, nil, "", refreshToken)
	if err != nil {
		return nil, err
	}
	_, err = uc.tokenManager.ValidateToken(u.RefreshToken())
	if err != nil {
		return nil, err
	}
	return u, nil
}

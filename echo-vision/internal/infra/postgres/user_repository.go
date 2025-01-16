package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/app/shared"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/user"
)

// FindUserByID implements repositories.UserRepository.
func (r *Repository) FindUserByID(
	ctx context.Context,
	tx repositories.Transaction,
	id uuid.UUID,
) (*domain.User, error) {
	c := r.resolveClient(tx)
	user, err := c.User.Query().
		Where(user.ID(id)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, shared.ErrUserNotFound
	}

	return userToDomain(user), err
}

// FindUserByEmail implements repositories.UserRepository.
func (r *Repository) FindUserByEmail(
	ctx context.Context,
	tx repositories.Transaction,
	email string,
) (*domain.User, error) {
	c := r.resolveClient(tx)
	user, err := c.User.Query().
		Where(user.Email(email)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, shared.ErrUserNotFound
	}

	return userToDomain(user), err
}

// SaveUser implements repositories.UserRepository.
func (r *Repository) SaveUser(
	ctx context.Context,
	tx repositories.Transaction,
	user *domain.User,
) (uuid.UUID, error) {
	c := r.resolveClient(tx)
	b := c.User.Create().
		SetID(user.ID()).
		SetFirstName(user.FirstName()).
		SetLastName(user.LastName()).
		SetEmail(user.Email())

	if user.Password() != "" {
		b.SetPassword(user.Password())
	}
	if user.AccessToken() != "" {
		b.SetAccessToken(user.AccessToken())
	}
	if user.RefreshToken() != "" {
		b.SetRefreshToken(user.RefreshToken())
	}

	u, err := b.Save(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return u.ID, nil
}

// UpdateTokens implements repositories.Repository.
func (r *Repository) UpdateTokens(
	ctx context.Context,
	tx repositories.Transaction,
	accessToken string,
	refreshToken string,
	userID uuid.UUID,
) error {
	c := r.resolveClient(tx)
	err := c.User.UpdateOneID(userID).
		SetAccessToken(accessToken).
		SetRefreshToken(refreshToken).
		Exec(ctx)
	return err
}

// UpdateUser implements repositories.Repository.
func (r *Repository) UpdateUser(ctx context.Context, tx repositories.Transaction, user *domain.User) error {
	c := r.resolveClient(tx)
	err := c.User.UpdateOneID(user.ID()).
		SetFirstName(user.FirstName()).
		SetLastName(user.LastName()).
		SetEmail(user.Email()).
		Exec(ctx)
	return err
}

// FindUserByTokens implements repositories.Repository.
func (r *Repository) FindUserByTokens(ctx context.Context, tx repositories.Transaction, accessToken string, refreshToken string) (*domain.User, error) {
	if accessToken == "" && refreshToken == "" {
		return nil, shared.ErrInvalidToken
	}

	c := r.resolveClient(tx)
	b := c.User.Query()

	if accessToken != "" {
		b.Where(user.AccessToken(accessToken))
	}
	if refreshToken != "" {
		b.Where(user.RefreshToken(refreshToken))
	}
	u, err := b.Only(ctx)
	if ent.IsNotFound(err) {
		return nil, shared.ErrUserNotFound
	}
	return userToDomain(u), err
}

// userToDomain transfer the ent object to the domain object
func userToDomain(entUser *ent.User) *domain.User {
	if entUser == nil {
		return nil
	}
	u := domain.NewUser(
		entUser.ID,
		entUser.FirstName,
		entUser.LastName,
		entUser.Email,
		entUser.CreatedAt,
		entUser.UpdatedAt,
	)
	u.SetHashedPassword(entUser.Password)
	u.SetTokens(entUser.AccessToken, entUser.RefreshToken)
	return u
}

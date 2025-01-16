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
		return nil, shared.ErrNotNound
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
		return nil, shared.ErrNotNound
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
	u, err := c.User.Create().
		SetID(user.ID()).
		SetFirstName(user.FirstName()).
		SetLastName(user.LastName()).
		SetEmail(user.Email()).
		Save(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return u.ID, nil
}

// userToDomain transfer the ent object to the domain object
func userToDomain(entUser *ent.User) *domain.User {
	if entUser == nil {
		return nil
	}
	return domain.NewUser(
		entUser.ID,
		entUser.FirstName,
		entUser.LastName,
		entUser.Email,
		entUser.CreatedAt,
		entUser.UpdatedAt,
	)
}

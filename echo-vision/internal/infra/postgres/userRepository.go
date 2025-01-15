package postgres

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/user"
)

// FindByEmail implements repositories.UserRepository.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.entClient.User.Query().
		Where(user.Email(email)).
		Only(ctx)
	return userToDomain(user), err
}

// Save implements repositories.UserRepository.
func (r *UserRepository) Save(ctx context.Context, tx repositories.Transaction, user *domain.User) error {
	err := r.entClient.User.Create().
		SetID(user.ID()).
		SetFirstName(user.FirstName()).
		SetLastName(user.LastName()).
		SetEmail(user.Email()).
		Exec(ctx)
	return err
}

func userToDomain(entUser *ent.User) *domain.User {
	return domain.NewUser(
		entUser.ID,
		entUser.FirstName,
		entUser.LastName,
		entUser.Email,
		entUser.CreatedAt,
		entUser.UpdatedAt,
	)
}

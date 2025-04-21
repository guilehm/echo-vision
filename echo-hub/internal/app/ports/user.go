package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
)

// UserPort defines the interface for user-related operations in the application layer.
// It provides methods for creating, saving, authenticating users, and managing tokens.
type UserPort interface {
	// CreateUser creates a new user with the provided details.
	CreateUser(ctx context.Context, firstName, lastName, email, password string) (*domain.User, error)

	// SaveUser persists the user entity and returns its UUID.
	SaveUser(ctx context.Context, user *domain.User) (uuid.UUID, error)

	// AuthenticateUser validates user credentials and returns the authenticated user.
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)

	// UserByAccessToken retrieves a user based on the provided access token.
	UserByAccessToken(ctx context.Context, accessToken string) (*domain.User, error)

	// UserByRefreshToken retrieves a user based on the provided refresh token.
	UserByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)

	// RefreshToken generates a new access token for the user using the refresh token.
	RefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
}

package postgres

import (
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
)

type UserRepository struct {
	entClient *ent.Client
}

type EventRepository struct {
	entClient *ent.Client
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(entClient *ent.Client) repositories.UserRepository {
	return &UserRepository{entClient: entClient}
}

// NewEventRepository creates a new instance of EventRepository.
func NewEventRepository(entClient *ent.Client) repositories.EventRepository {
	return &EventRepository{entClient: entClient}
}

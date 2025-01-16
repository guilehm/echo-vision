package postgres

import (
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
)

type Repository struct {
	entClient *ent.Client
}

// NewRepository creates a new instance of Repository.
func NewRepository(entClient *ent.Client) repositories.Repository {
	return &Repository{entClient: entClient}
}

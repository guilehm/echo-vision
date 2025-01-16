package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/repositories"
)

type ManageEvents struct {
	Repository repositories.Repository
}

func (uc *ManageEvents) FindEventByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Event, error) {
	return uc.Repository.FindEventByID(ctx, nil, id)
}

func (uc *ManageEvents) SaveEvent(
	ctx context.Context,
	event *domain.Event,
) (uuid.UUID, error) {
	return uc.Repository.SaveEvent(ctx, nil, event)
}

func NewManageEventsUseCase(repository repositories.Repository) ports.EventPort {
	return &ManageEvents{
		Repository: repository,
	}
}

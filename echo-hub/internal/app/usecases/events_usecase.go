package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
)

type ManageEvents struct {
	Repository repositories.Repository
	publisher  messaging.Publisher
}

func NewManageEventsUseCase(repository repositories.Repository, publisher messaging.Publisher) ports.EventPort {
	return &ManageEvents{
		Repository: repository,
		publisher:  publisher,
	}
}

// CreateEvent implements ports.EventPort.
func (uc *ManageEvents) CreateEvent(
	ctx context.Context,
	userID uuid.UUID,
	eventType string,
	subType string,
) (*domain.Event, error) {
	now := time.Now()
	event := domain.NewEvent(
		userID,
		uuid.New(),
		domain.EventType(eventType),
		domain.EventSubType(subType),
		nil,
		nil,
		domain.EventStatusPending,
		now,
		now,
	)
	if err := event.Validate(); err != nil {
		return nil, err
	}
	return event, nil
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

// EventsByUser implements ports.EventPort.
func (uc *ManageEvents) EventsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	return uc.Repository.FindEventsByUserID(ctx, nil, userID)
}

package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	"github.com/rotisserie/eris"
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
	file *valueobjects.File,
) (*domain.Event, error) {
	now := time.Now()
	event := domain.NewEvent(
		userID,
		uuid.New(),
		hubevents.EventType(eventType),
		hubevents.EventSubType(subType),
		nil,
		hubevents.EventStatusPending,
		file,
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
	id, err := uc.Repository.SaveEvent(ctx, nil, event)
	if err != nil {
		return id, err
	}

	payload, err := ports.MapEventToMessage(event)
	if err != nil {
		return id, eris.Wrap(err, "failed to map event to json message")
	}

	// TODO: only publish this message on commit
	err = uc.publisher.Publish(ctx, messaging.Message{
		Topic:   hubevents.BuildEventCreatedTopic(event.EventType()),
		Payload: payload,
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

// EventsByUser implements ports.EventPort.
func (uc *ManageEvents) EventsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	return uc.Repository.FindEventsByUserID(ctx, nil, userID)
}

package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	"github.com/rotisserie/eris"

	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
)

type ManageEvents struct {
	Repository  repositories.Repository
	publisher   ports.PublisherPort
	fileStorage filestorage.FileStoragePort
}

func NewManageEventsUseCase(
	repository repositories.Repository,
	publisher ports.PublisherPort,
	fs filestorage.FileStoragePort,
) ports.EventPort {
	return &ManageEvents{
		Repository:  repository,
		publisher:   publisher,
		fileStorage: fs,
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
	event, err := uc.Repository.FindEventByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}

	if event.File() == nil {
		return event, nil
	}

	url, err := uc.fileStorage.GenerateFileURL(event.File().FileKey())
	if err != nil {
		return nil, eris.Wrap(err, "failed to generate file URL")
	}
	event.SetFileURL(url)

	return event, nil
}

func (uc *ManageEvents) SaveEvent(
	ctx context.Context,
	event *domain.Event,
) (uuid.UUID, error) {
	id, err := uc.Repository.SaveEvent(ctx, nil, event)
	if err != nil {
		return id, err
	}

	err = uc.publisher.PublishEventCreated(ctx, event)
	if err != nil {
		return id, err
	}
	return id, nil
}

// EventsByUser implements ports.EventPort.
func (uc *ManageEvents) EventsByUser(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*domain.Event, string, error) {
	return uc.Repository.FindEventsByUserID(ctx, nil, userID, limit, cursor)
}

// HandleImageAnalysisStatusUpdate implements ports.EventPort.
func (uc *ManageEvents) HandleImageAnalysisStatusUpdate(
	ctx context.Context,
	id uuid.UUID,
	status hubevents.EventStatus,
	result json.RawMessage,
) error {
	event, err := uc.Repository.FindEventByID(ctx, nil, id)
	if err != nil {
		return eris.Wrap(err, "failed to find event by ID")
	}

	// set event status
	if err := event.SetStatus(status); err != nil {
		return eris.Wrap(err, "failed to set event status")
	}

	// set result
	event.SetResult(result)

	// validate event before saving
	if err := event.Validate(); err != nil {
		return eris.Wrap(err, "failed to validate event")
	}

	// persist the event
	if err := uc.Repository.UpdateEvent(ctx, nil, event); err != nil {
		return eris.Wrap(err, "failed to save event")
	}

	// TODO: only publish this message on commit
	err = uc.publisher.PublishEventStatusUpdated(ctx, event, result)
	if err != nil {
		return eris.Wrap(err, "failed to publish event status update")
	}
	return nil
}

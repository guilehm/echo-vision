package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type EventPort interface {
	FindEventByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	CreateEvent(ctx context.Context, userID uuid.UUID, eventType, subType string, file *valueobjects.File) (*domain.Event, error)
	SaveEvent(ctx context.Context, event *domain.Event) (uuid.UUID, error)
	EventsByUser(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]*domain.Event, string, error)
	HandleImageAnalysisStatusUpdate(ctx context.Context, id uuid.UUID, status hubevents.EventStatus, result json.RawMessage) error
}

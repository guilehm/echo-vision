package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type EventPort interface {
	FindEventByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	CreateEvent(ctx context.Context, userID uuid.UUID, eventType, subType string) (*domain.Event, error)
	SaveEvent(ctx context.Context, event *domain.Event) (uuid.UUID, error)
	EventsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error)
}

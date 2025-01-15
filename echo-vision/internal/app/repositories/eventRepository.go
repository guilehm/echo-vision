package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type EventRepository interface {
	SaveEvent(ctx context.Context, tx Transaction, event *domain.Event) error
	FindEventByID(ctx context.Context, tx Transaction, id uuid.UUID) (*domain.Event, error)
}

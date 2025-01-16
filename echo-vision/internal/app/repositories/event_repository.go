package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type EventRepository interface {
	SaveEvent(ctx context.Context, tx Transaction, event *domain.Event) (uuid.UUID, error)
	FindEventByID(ctx context.Context, tx Transaction, id uuid.UUID) (*domain.Event, error)
	FindEventsByUserID(
		ctx context.Context,
		tx Transaction,
		userID uuid.UUID,
	) ([]*domain.Event, error)
}

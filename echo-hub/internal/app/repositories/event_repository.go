package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type EventRepository interface {
	SaveEvent(ctx context.Context, tx Transaction, event *domain.Event) (uuid.UUID, error)
	FindEventByID(ctx context.Context, tx Transaction, id uuid.UUID) (*domain.Event, error)
	FindEventsByUserID(
		ctx context.Context,
		tx Transaction,
		userID uuid.UUID,
	) ([]*domain.Event, error)
	UpdateEventStatus(
		ctx context.Context,
		tx Transaction,
		id uuid.UUID,
		status hubevents.EventStatus,
	) error
	UpdateEvent(
		ctx context.Context,
		tx Transaction,
		event *domain.Event,
	) error
}

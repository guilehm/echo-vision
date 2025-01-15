package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/logging"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/event"
)

var logger = logging.NewLogger()

// Save implements repositories.EventRepository.
func (r *Repository) SaveEvent(ctx context.Context, tx repositories.Transaction, e *domain.Event) error {
	err := r.entClient.Event.Create().
		SetID(e.ID()).
		SetType(event.Type(e.EventType())).
		SetSubType(event.SubType(e.SubType())).
		SetStatus(event.Status(e.Status())).
		SetPayload(e.Payload()).
		SetResult(e.Result()).
		SetCreatedAt(e.CreatedAt()).
		SetUpdatedAt(e.UpdatedAt()).
		Exec(ctx)
	return err
}

// FindEventByID implements repositories.EventRepository.
func (r *Repository) FindEventByID(ctx context.Context, tx repositories.Transaction, id uuid.UUID) (*domain.Event, error) {
	e, err := r.entClient.Event.Query().
		Where(event.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return eventToDomain(e), nil
}

func eventToDomain(e *ent.Event) *domain.Event {
	return domain.NewEvent(
		e.ID,
		domain.EventType(e.Type.String()),
		domain.EventSubType(e.SubType.String()),
		e.Payload,
		e.Result,
		domain.EventStatus(e.Status.String()),
		e.CreatedAt,
		e.UpdatedAt,
	)
}

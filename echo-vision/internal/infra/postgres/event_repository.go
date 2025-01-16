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
func (r *Repository) SaveEvent(ctx context.Context, tx repositories.Transaction, e *domain.Event) (uuid.UUID, error) {
	entEvent, err := r.entClient.Event.Create().
		SetUserID(e.User().ID()).
		SetID(e.ID()).
		SetType(event.Type(e.EventType())).
		SetSubType(event.SubType(e.SubType())).
		SetStatus(event.Status(e.Status())).
		SetPayload(e.Payload()).
		SetResult(e.Result()).
		SetCreatedAt(e.CreatedAt()).
		SetUpdatedAt(e.UpdatedAt()).
		Save(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return entEvent.ID, nil
}

// FindEventByID implements repositories.EventRepository.
func (r *Repository) FindEventByID(ctx context.Context, tx repositories.Transaction, id uuid.UUID) (*domain.Event, error) {
	e, err := r.entClient.Event.Query().
		Where(event.ID(id)).
		Only(ctx)
	return eventToDomain(e), err
}

func eventToDomain(e *ent.Event) *domain.Event {
	if e == nil {
		return nil
	}
	return domain.NewEvent(
		userToDomain(e.Edges.User),
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

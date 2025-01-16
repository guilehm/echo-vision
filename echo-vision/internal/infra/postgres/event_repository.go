package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/app/shared"
	"github.com/guilehm/echo-vision/internal/infra/logging"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/event"
)

var logger = logging.NewLogger()

// Save implements repositories.EventRepository.
func (r *Repository) SaveEvent(
	ctx context.Context,
	tx repositories.Transaction,
	e *domain.Event,
) (uuid.UUID, error) {
	c := r.resolveClient(tx)
	entEvent, err := c.Event.Create().
		SetUserID(e.UserID()).
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
func (r *Repository) FindEventByID(
	ctx context.Context,
	tx repositories.Transaction,
	id uuid.UUID,
) (*domain.Event, error) {
	c := r.resolveClient(tx)
	e, err := c.Event.Query().
		Where(event.ID(id)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, shared.ErrNotFound
	}
	return eventToDomain(e), err
}

// FindEventsByUserID implements repositories.EventRepository.
func (r *Repository) FindEventsByUserID(
	ctx context.Context,
	tx repositories.Transaction,
	userID uuid.UUID,
) ([]*domain.Event, error) {
	// TODO: paginate with cursor
	c := r.resolveClient(tx)
	events, err := c.Event.Query().
		Where(event.UserID(userID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var result []*domain.Event
	for _, e := range events {
		result = append(result, eventToDomain(e))
	}
	return result, nil
}

// eventToDomain transfer the ent object to the domain object
func eventToDomain(e *ent.Event) *domain.Event {
	if e == nil {
		return nil
	}
	return domain.NewEvent(
		e.UserID,
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

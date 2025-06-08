package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/event"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

var logger = logging.NewLogger()

// Save implements repositories.EventRepository.
func (r *Repository) SaveEvent(
	ctx context.Context,
	tx repositories.Transaction,
	e *domain.Event,
) (uuid.UUID, error) {
	c := r.resolveClient(tx)

	var err error
	var file *ent.File

	if e.File() != nil {
		file, err = c.File.Create().
			SetID(uuid.New()).
			SetFilename(e.File().Filename).
			SetFilepath(e.File().Filepath).
			SetFilesize(e.File().Filesize).
			SetContentType(e.File().ContentType).
			Save(ctx)
		if err != nil {
			return uuid.Nil, err
		}
	}

	builder := c.Event.Create().
		SetUserID(e.UserID()).
		SetID(e.ID()).
		SetType(event.Type(e.EventType())).
		SetSubType(event.SubType(e.SubType())).
		SetStatus(event.Status(e.Status())).
		SetResult(e.Result()).
		SetCreatedAt(e.CreatedAt()).
		SetUpdatedAt(e.UpdatedAt())

	if file != nil {
		builder.SetFile(file)
	}

	entEvent, err := builder.Save(ctx)
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
		WithFile().
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
	limit int,
	cursor string,
) ([]*domain.Event, string, error) {
	c := r.resolveClient(tx)

	var err error

	// build the default query
	query := c.Event.Query().
		Where(event.UserID(userID)).
		WithFile()

	// apply limit and cursor filters if provided
	query, err = eventsCursorFilter(query, cursor)
	if err != nil {
		return nil, "", err
	}

	// execute the query
	events, err := query.
		Order(ent.Desc(event.FieldCreatedAt), ent.Desc(event.FieldID)).
		Limit(limit + 1).
		All(ctx)
	if err != nil {
		return nil, "", err
	}

	// map results to domain
	results, nextCursor := eventResultsAndCursor(events, limit)
	return results, nextCursor, nil
}

// eventToDomain transfer the ent object to the domain object
func eventToDomain(e *ent.Event) *domain.Event {
	if e == nil {
		return nil
	}
	var file *valueobjects.File
	if e.Edges.File != nil {
		file = valueobjects.NewFile(
			e.Edges.File.Filepath,
			e.Edges.File.Filename,
			e.Edges.File.ContentType,
			e.Edges.File.Filesize,
		)
	}
	return domain.NewEvent(
		e.UserID,
		e.ID,
		hubevents.EventType(e.Type.String()),
		hubevents.EventSubType(e.SubType.String()),
		e.Result,
		hubevents.EventStatus(e.Status.String()),
		file,
		e.CreatedAt,
		e.UpdatedAt,
	)
}

func (r *Repository) UpdateEventStatus(
	ctx context.Context,
	tx repositories.Transaction,
	id uuid.UUID,
	status hubevents.EventStatus,
) error {
	c := r.resolveClient(tx)
	return c.Event.UpdateOneID(id).
		SetStatus(event.Status(status)).
		Exec(ctx)
}

func (r *Repository) UpdateEvent(
	ctx context.Context,
	tx repositories.Transaction,
	e *domain.Event,
) error {
	c := r.resolveClient(tx)
	return c.Event.UpdateOneID(e.ID()).
		SetType(event.Type(e.EventType())).
		SetSubType(event.SubType(e.SubType())).
		SetStatus(event.Status(e.Status())).
		SetResult(e.Result()).
		Exec(ctx)
}

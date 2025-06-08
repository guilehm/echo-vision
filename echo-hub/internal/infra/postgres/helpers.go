package postgres

import (
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/event"
)

func eventsCursorFilter(query *ent.EventQuery, cursor string) (*ent.EventQuery, error) {
	if cursor == "" {
		return query, nil
	}

	// decode cursor: createdAt|id
	createdAt, id, err := decodeCursor(cursor)
	if err != nil {
		return query, shared.ErrInvalidCursor
	}
	return query.Where(
		event.Or(
			event.And(
				event.CreatedAtLT(createdAt),
			),
			event.And(
				event.IDLT(id),
				event.CreatedAtEQ(createdAt),
			),
		),
	), nil
}

func eventResultsAndCursor(events []*ent.Event, limit int) ([]*domain.Event, string) {
	var result []*domain.Event
	var nextCursor string
	for i, e := range events {
		result = append(result, eventToDomain(e))
		if i+1 == limit && len(events) > limit {
			// if this is the last event, set the next cursor
			nextCursor = encodeCursor(e.CreatedAt, e.ID)
			break
		}
	}
	return result, nextCursor
}

package postgres

import (
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent/event"
)

func applyCursorFilter(query *ent.EventQuery, cursor string) (*ent.EventQuery, error) {
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

package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

type EventPort interface {
	FindEventByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	SaveEvent(ctx context.Context, event *domain.Event) (uuid.UUID, error)
}

type EventCreateInput struct {
	EventType string `json:"eventType"`
	SubType   string `json:"subType"`
}

type EventResponse struct {
	UserID    uuid.UUID `json:"userID"`
	ID        uuid.UUID `json:"id"`
	EventType string    `json:"eventType"`
	SubType   string    `json:"subType"`
	Status    string    `json:"status"`
}

type EventListResponse struct {
	Events []*EventResponse `json:"events"`
	Page   int              `json:"page"`
	Count  int              `json:"count"`
}

// type EventPort interface {
// 	FindEventByID(ctx context.Context, id uuid.UUID) (*EventResponse, error)
// 	SaveEvent(ctx context.Context, event *domain.Event) (*EventResponse, error)
// 	ListEvents(ctx context.Context) (*EventListResponse, error)
// }

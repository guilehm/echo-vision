package ports

import (
	"context"
	"encoding/json"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
)

type PublisherPort interface {
	PublishEventCreated(ctx context.Context, event *domain.Event) error
	PublishEventStatusUpdated(ctx context.Context, event *domain.Event, data json.RawMessage) error
}

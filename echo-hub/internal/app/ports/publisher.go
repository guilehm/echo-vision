package ports

import (
	"context"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
)

type PublisherPort interface {
	PublishEventCreated(ctx context.Context, event *domain.Event) error
}

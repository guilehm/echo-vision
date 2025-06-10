package publishers

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	"github.com/rotisserie/eris"
)

// PulbishEventCreated implements ports.PublisherPort.
func (p *PublisherGroup) PublishEventCreated(ctx context.Context, event *domain.Event) error {
	payload, err := ports.MapEventToMessage(event)
	if err != nil {
		return eris.Wrap(err, "failed to map event to json message")
	}
	return p.publisher.Publish(ctx, messaging.Message{
		Topic:   hubevents.BuildEventCreatedTopic(event.EventType()),
		Payload: payload,
	})
}

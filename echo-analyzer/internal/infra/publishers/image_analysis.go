package publishers

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
	"github.com/rotisserie/eris"
)

// PublishImageAnalysisStatusUpdate implements ports.PublisherPort.
func (p *PublisherGroup) PublishImageAnalysisStatusUpdate(ctx context.Context, eventID uuid.UUID, message hubevents.EventStatusUpdateMessage) error {
	payload, err := message.ToJSON()
	if err != nil {
		return eris.Wrap(err, "failed to convert message to JSON")
	}

	err = p.publisher.Publish(ctx, messaging.Message{
		Topic:   hubevents.BuildEventStatusUpdatedTopic(message.Type, message.Status),
		Payload: payload,
	})
	if err != nil {
		return eris.Wrap(err, "failed to publish message for event")
	}
	return nil
}

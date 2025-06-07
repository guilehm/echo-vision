package publishers

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/rotisserie/eris"
)

// PublishImageAnalysisStatusUpdate implements ports.PublisherPort.
func (p *PublisherGroup) PublishImageAnalysisStatusUpdate(
	ctx context.Context,
	eventID uuid.UUID,
	status analyzerevents.EventStatus,
	data json.RawMessage,
) error {
	message := analyzerevents.EventImageAnalysisStatusUpdateMessage{
		ID:     eventID,
		Status: status,
		Data:   data,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return eris.Wrap(err, "failed to marshal event message")
	}
	err = p.publisher.Publish(ctx, messaging.Message{
		Topic:   analyzerevents.BuildEventImageAnalysisStatusUpdatedTopic(status),
		Payload: payload,
	})
	if err != nil {
		return eris.Wrap(err, "failed to publish message for event")
	}
	return nil
}

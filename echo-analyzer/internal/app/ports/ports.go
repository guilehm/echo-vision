package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type ConsumerPort interface {
	ProcessImageAnalysis(topic string, message hubevents.EventMessage) messaging.HandlerResponse
}

type PublisherPort interface {
	PublishImageAnalysisStatusUpdate(ctx context.Context, eventID uuid.UUID, message hubevents.EventStatusUpdateMessage) error
}

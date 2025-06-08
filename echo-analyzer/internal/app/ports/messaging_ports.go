package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type ConsumerPort interface {
	ProcessImageAnalysis(topic string, message hubevents.EventMessage) messaging.HandlerResponse
}

type PublisherPort interface {
	PublishImageAnalysisStatusUpdate(ctx context.Context, eventID uuid.UUID, status analyzerevents.EventStatus, data json.RawMessage) error
}

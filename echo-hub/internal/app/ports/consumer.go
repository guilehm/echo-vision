package ports

import (
	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

type ConsumerPort interface {
	ImageAnalysisStatusUpdate(topic string, id uuid.UUID, status string) messaging.HandlerResponse
}

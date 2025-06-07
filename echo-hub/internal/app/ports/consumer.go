package ports

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

type ConsumerPort interface {
	ImageAnalysisStatusUpdate(id uuid.UUID, status string, result json.RawMessage) messaging.HandlerResponse
}

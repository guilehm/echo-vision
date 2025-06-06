package publishers

import (
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

var logger = logging.NewLogger()

type PublisherGroup struct {
	irs       ports.ImageRecognitionServicePort
	publisher messaging.Publisher
}

func NewPublisherGroup(irs ports.ImageRecognitionServicePort, publisher messaging.Publisher) ports.PublisherPort {
	return &PublisherGroup{
		irs:       irs,
		publisher: publisher,
	}
}

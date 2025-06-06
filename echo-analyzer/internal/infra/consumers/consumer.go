package consumers

import (
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/logging"
)

var logger = logging.NewLogger()

type ConsumerGroup struct {
	irs       ports.ImageRecognitionServicePort
	publisher ports.PublisherPort
}

func NewConsumerGroup(irs ports.ImageRecognitionServicePort, publisher ports.PublisherPort) ports.ConsumerPort {
	return &ConsumerGroup{
		irs:       irs,
		publisher: publisher,
	}
}

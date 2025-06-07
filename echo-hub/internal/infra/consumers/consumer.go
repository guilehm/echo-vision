package consumers

import (
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

var logger = logging.NewLogger()

type ConsumerGroup struct {
	EventUseCase ports.EventPort
}

func NewConsumerGroup(eventUseCase ports.EventPort) ports.ConsumerPort {
	return &ConsumerGroup{
		EventUseCase: eventUseCase,
	}
}

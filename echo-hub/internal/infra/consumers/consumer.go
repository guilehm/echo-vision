package consumers

import (
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

var logger = logging.NewLogger()

type ConsumerGroup struct{}

func NewConsumerGroup() ports.ConsumerPort {
	return &ConsumerGroup{}
}

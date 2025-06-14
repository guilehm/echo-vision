package consumers

import (
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/pkg/logging"
)

var logger = logging.NewLogger()

type ConsumerGroup struct {
	publisher            ports.PublisherPort
	imageAnalysisUseCase ports.ImageAnalysisPort
}

func NewConsumerGroup(imageAnalysisUseCase ports.ImageAnalysisPort, publisher ports.PublisherPort) ports.ConsumerPort {
	return &ConsumerGroup{
		publisher:            publisher,
		imageAnalysisUseCase: imageAnalysisUseCase,
	}
}

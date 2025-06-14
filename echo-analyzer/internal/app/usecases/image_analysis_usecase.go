package usecases

import (
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/pkg/logging"
)

var logger = logging.NewLogger()

type ImageAnalysisUseCase struct {
	publisher ports.PublisherPort
	irs       ports.ImageRecognitionServicePort
}

func NewImageAnalysisUseCase(publisher ports.PublisherPort, irs ports.ImageRecognitionServicePort) ports.ImageAnalysisPort {
	return ImageAnalysisUseCase{
		publisher: publisher,
		irs:       irs,
	}
}

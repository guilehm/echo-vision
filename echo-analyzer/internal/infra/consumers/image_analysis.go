package consumers

import (
	"context"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// ProcessImageAnalysis implements ports.ConsumerPort.
func (c *ConsumerGroup) ProcessImageAnalysis(topic string, message hubevents.EventMessage) messaging.HandlerResponse {
	ctx := context.Background()
	if message.File == nil {
		logger.Error("message does not contain a file")
		return messaging.DeadLetter
	}

	if message.Type != hubevents.EventTypeImageAnalysis {
		logger.Error("message type is not image analysis", slog.String("type", string(message.Type)))
		return messaging.DeadLetter
	}

	switch message.SubType {
	case hubevents.EventSubTypeDetectLabels:
		_, err := c.imageAnalysisUseCase.DetectLabels(ctx, message.ID, message.File.Filepath)
		if err != nil {
			logger.Error("could not detect labels", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
	case hubevents.EventSubTypeDetectFaces:
		_, err := c.imageAnalysisUseCase.DetectFaces(ctx, message.ID, message.File.Filepath)
		if err != nil {
			logger.Error("could not detect faces", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
	}

	return messaging.Success
}

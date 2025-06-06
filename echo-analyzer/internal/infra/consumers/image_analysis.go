package consumers

import (
	"context"
	"encoding/json"
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

	var data json.RawMessage
	switch message.SubType {
	case hubevents.EventSubTypeDetectLabels:
		labels, err := c.irs.DetectLabels(message.File.Filepath)
		if err != nil {
			logger.Error("could not detect labels: ", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
		data = toRawMessage(labels)

	case hubevents.EventSubTypeDetectFaces:
		faces, err := c.irs.DetectFaces(message.File.Filepath)
		if err != nil {
			logger.Error("could not detect faces ", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
		data = toRawMessage(faces)
	}

	err := c.publisher.PublishImageAnalysisStatusUpdate(ctx, message.ID, hubevents.EventStatusUpdateMessage{
		ID:     message.ID,
		Type:   message.Type,
		Status: hubevents.EventStatusCompleted,
		Data:   data,
	})
	if err != nil {
		logger.Error("could not publish image analysis status update", slog.String("error", err.Error()))
		return messaging.DeadLetter
	}
	return messaging.Success
}

func toRawMessage(data any) json.RawMessage {
	raw, err := json.Marshal(data)
	if err != nil {
		logger.Error("could not marshal data to raw message", slog.String("error", err.Error()))
		return nil
	}
	return raw
}

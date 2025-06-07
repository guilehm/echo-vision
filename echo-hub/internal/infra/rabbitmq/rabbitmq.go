package rabbitmqadapter

import (
	"context"
	"encoding/json"
	"log/slog"

	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

var logger = logging.NewLogger()

type RabbitMQAdapter struct {
	consumer ports.ConsumerPort
}

func (r *RabbitMQAdapter) Topics() []string {
	return []string{
		analyzerevents.EventImageAnalysisStatusUpdatedProcessing,
		analyzerevents.EventImageAnalysisStatusUpdatedCompleted,
		analyzerevents.EventImageAnalysisStatusUpdatedFailed,
	}
}

func (r *RabbitMQAdapter) Handle(ctx context.Context, msg messaging.Message) messaging.HandlerResponse {
	switch msg.Topic {
	case
		analyzerevents.EventImageAnalysisStatusUpdatedProcessing,
		analyzerevents.EventImageAnalysisStatusUpdatedCompleted,
		analyzerevents.EventImageAnalysisStatusUpdatedFailed:
		var message analyzerevents.EventImageAnalysisStatusUpdateMessage
		if err := json.Unmarshal(msg.Payload, &message); err != nil {
			logger.Error("could not unmarshal event", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
		logger.Info("received image analysis status update",
			slog.String("id", message.ID.String()),
			slog.String("status", string(message.Status)))
		return r.consumer.ImageAnalysisStatusUpdate(
			message.ID,
			message.Status.String(),
			message.Data,
		)
	default:
		return messaging.DeadLetter
	}
}

func NewRabbitMQAdapter(consumer ports.ConsumerPort) *RabbitMQAdapter {
	return &RabbitMQAdapter{
		consumer: consumer,
	}
}

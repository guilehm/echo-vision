package rabbitmqadapter

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

var logger = logging.NewLogger()

type RabbitMQAdapter struct {
	consumer ports.ConsumerPort
}

func (r *RabbitMQAdapter) Topics() []string {
	return []string{
		hubevents.EventImageAnalysisStatusUpdatedGeneric,
	}
}

func (r *RabbitMQAdapter) Handle(ctx context.Context, msg messaging.Message) messaging.HandlerResponse {
	switch msg.Topic {
	case
		hubevents.EventImageAnalysisStatusUpdatedCompleted,
		hubevents.EventImageAnalysisStatusUpdatedFailed,
		hubevents.EventImageAnalysisStatusUpdatedProcessing,
		hubevents.EventImageAnalysisStatusUpdatedPending:
		var message hubevents.EventStatusUpdateMessage
		if err := json.Unmarshal(msg.Payload, &message); err != nil {
			logger.Error("could not unmarshal event", slog.String("error", err.Error()))
			return messaging.DeadLetter
		}
		logger.Info("received image analysis status update",
			slog.String("id", message.ID.String()),
			slog.String("status", string(message.Status)))
		return r.consumer.ImageAnalysisStatusUpdate(msg.Topic, message)
	default:
		return messaging.DeadLetter
	}
}

func NewRabbitMQAdapter(consumer ports.ConsumerPort) *RabbitMQAdapter {
	return &RabbitMQAdapter{
		consumer: consumer,
	}
}

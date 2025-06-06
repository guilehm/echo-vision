package rabbitmqadapter

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"

	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type RabbitMQAdapter struct {
	consumer  ports.ConsumerPort
	publisher ports.PublisherPort
}

func (r *RabbitMQAdapter) Topics() []string {
	return []string{
		hubevents.EventImageAnalysCreated,
	}
}

func (r *RabbitMQAdapter) Handle(ctx context.Context, msg messaging.Message) messaging.HandlerResponse {
	logger := logging.NewLogger().With(slog.String("topic", msg.Topic))

	switch msg.Topic {
	case hubevents.EventImageAnalysCreated:
		var message hubevents.EventMessage
		if err := json.Unmarshal(msg.Payload, &message); err != nil {
			logger.Error(
				"could not unmarshal event",
				slog.String("error", err.Error()),
			)
			return messaging.DeadLetter
		}
		return r.consumer.ProcessImageAnalysis(msg.Topic, message)
	default:
		return messaging.DeadLetter
	}
}

func NewRabbitMQAdapter(consumer ports.ConsumerPort, publisher ports.PublisherPort) *RabbitMQAdapter {
	return &RabbitMQAdapter{
		consumer:  consumer,
		publisher: publisher,
	}
}

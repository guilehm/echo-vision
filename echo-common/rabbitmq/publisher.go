package rabbitmq

import (
	"context"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	config     *Config
	controller *ChannelController

	unpublishedMessages chan messaging.Message
}

// StartPublisher implements Publisher.
func (r *RabbitMQPublisher) StartPublisher(ctx context.Context) error {
	go r.publishLoop()
	go r.controller.reconnectLoop()

	r.controller.Resume()
	return nil
}

// Close implements Publisher.
func (r *RabbitMQPublisher) Close() error {
	return r.controller.ch.Close()
}

// Publish implements Publisher.
func (r *RabbitMQPublisher) Publish(ctx context.Context, message messaging.Message) error {
	r.unpublishedMessages <- message
	return nil
}

func (r *RabbitMQPublisher) publish(message messaging.Message) {
	logger := r.config.Logger
	ctx, cancel := context.WithTimeout(context.Background(), r.config.PublishTimeout)

	// publish message
	deferredConfirmation, err := r.controller.ch.PublishWithDeferredConfirmWithContext(
		ctx,
		r.config.ExchangeName.String(),
		message.Topic,
		true,
		false,
		amqp091.Publishing{
			Headers:      amqpTableFromMap(message.Headers),
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         message.Payload,
		},
	)
	if err != nil {
		// the connection is down, we should pause the publisher and republish the message
		logger.Error("could not publish message",
			slog.String("error", err.Error()),
			slog.String("topic", message.Topic),
		)

		// republish message
		r.controller.Pause()
		r.unpublishedMessages <- message
		cancel()
		return
	}

	// wait for confirmation
	confirmed, err := deferredConfirmation.WaitContext(ctx)
	cancel()
	if err != nil {
		logger.Error("could not confirm message",
			slog.String("error", err.Error()),
			slog.String("topic", message.Topic),
		)

		// republish message
		r.unpublishedMessages <- message
	}
	if !confirmed {
		logger.Error("message was not confirmed, republishing",
			slog.String("topic", message.Topic),
		)

		// republish message
		r.unpublishedMessages <- message
	}
	logger.Info("message published", slog.String("topic", message.Topic))
}

func (r *RabbitMQPublisher) publishLoop() {
	for {
		r.controller.WaitForResume()
		go r.publish(<-r.unpublishedMessages)
	}
}

func amqpTableFromMap(m map[string]string) amqp091.Table {
	table := make(amqp091.Table, len(m))
	for k, v := range m {
		table[k] = v
	}
	return table
}

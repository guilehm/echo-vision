package rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	config Config

	ch                  *amqp091.Channel
	unpublishedMessages chan Message
	publisherFunc       func(Message)
}

// Close implements Publisher.
func (r *RabbitMQPublisher) Close() error {
	panic("unimplemented")
}

// Publish implements Publisher.
func (r *RabbitMQPublisher) Publish(ctx context.Context, topic string, message Message) error {
	panic("unimplemented")
}

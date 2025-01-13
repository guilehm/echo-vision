package rabbitmq

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQChannel struct {
	channel *amqp.Channel
}

func (rc *RabbitMQChannel) Publish(ctx context.Context, topic string, message ports.Message) error {
	panic("unimplemented")
}

func (rc *RabbitMQChannel) Subscribe(ctx context.Context, topic string, handler func(msg ports.Message) error) error {
	panic("unimplemented")
}

func (rc *RabbitMQChannel) Close() error {
	return rc.channel.Close()
}

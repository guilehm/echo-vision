package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQChannel struct {
	channel *amqp.Channel
}

func (rc *RabbitMQChannel) Publish(ctx context.Context, topic string, message Message) error {
	body, err := json.Marshal(message.Payload)
	if err != nil {
		return ErrCouldNotDecodeMessage
	}
}

func (rc *RabbitMQChannel) Subscribe(ctx context.Context, topic string, handler func(msg Message) error) error {
	panic("unimplemented")
}

func (rc *RabbitMQChannel) Close() error {
	return rc.channel.Close()
}

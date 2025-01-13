package rabbitmq

import (
	"errors"
	"os"

	"github.com/guilehm/echo-vision/internal/app/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeType string

func (et ExchangeType) String() string {
	return string(et)
}

const (
	ExchangeTypeDirect  ExchangeType = "direct"
	ExchangeTypeFanout  ExchangeType = "fanout"
	ExchangeTypeTopic   ExchangeType = "topic"
	ExchangeTypeHeaders ExchangeType = "headers"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
}

func NewRabbitMQAdapter() (ports.AsyncMessagingPort, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		return nil, errors.New("RABBITMQ_URL is required")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}
	return &RabbitMQAdapter{connection: conn}, nil
}

func (r *RabbitMQAdapter) CreateChannel() (ports.MessagingChannel, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQChannel{channel: ch}, nil
}

func (r *RabbitMQAdapter) Close() error {
	return r.connection.Close()
}

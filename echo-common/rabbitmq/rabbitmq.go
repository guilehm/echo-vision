package rabbitmq

import (
	"errors"
	"os"

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

type RabbitMQClient struct {
	connection *amqp.Connection
}

func NewRabbitMQClient() (AsyncMessagingPort, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		return nil, errors.New("RABBITMQ_URL is required")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{connection: conn}, nil
}

func (r *RabbitMQClient) CreateChannel() (MessagingChannel, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQChannel{channel: ch}, nil
}

func (r *RabbitMQClient) Close() error {
	return r.connection.Close()
}

package rabbitmq

import (
	"errors"

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

type Config struct {
	ExchangeName        ExchangeName
	QueueName           QueueName
	ConsumerName        ConsumerName
	URL                 string
	ConcurrentConsumers int
	PrefetchCount       int // PrefetchCount is the number of messages to fetch from the queue at a time.
}

type RabbitMQClient struct {
	connection *amqp.Connection
	config     Config
}

func (r *RabbitMQClient) NewPublisher() (Publisher, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQPublisher{
		ch: ch,
	}, nil
}

func (r *RabbitMQClient) NewConsumer() (Publisher, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQConsumer{
		ch:     ch,
		config: r.config,
	}, nil
}

// CreateConsumer implements AsyncMessagingPort.
func (r *RabbitMQClient) CreateConsumer() (Consumer, error) {
	panic("unimplemented")
}

// CreatePublisher implements AsyncMessagingPort.
func (r *RabbitMQClient) CreatePublisher() (Publisher, error) {
	panic("unimplemented")
}

func NewRabbitMQClient(config Config) (AsyncMessagingPort, error) {
	if config.URL == "" {
		return nil, errors.New("RabbitMQ URL is required")
	}

	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{
		connection: conn,
		config:     config,
	}, nil
}

// func (r *RabbitMQClient) CreateChannel() (MessagingChannel, error) {
// 	ch, err := r.connection.Channel()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &RabbitMQChannel{channel: ch}, nil
// }

func (r *RabbitMQClient) Close() error {
	return r.connection.Close()
}

package rabbitmq

import (
	"log/slog"
	"sync"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rotisserie/eris"
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
	config *Config
}

func NewRabbitMQClient(url string, logger *slog.Logger, opts ...ConfigOpt) (AsyncMessagingPort, error) {
	config := newRabbitMQConfig(url, logger, opts...)
	if config.URL == "" {
		return nil, ErrRabbitMQURLIsRequired
	}

	if config.Logger == nil {
		return nil, ErrorLoggerIsRequired
	}
	return &RabbitMQClient{
		config: config,
	}, nil
}

func createConnection(config *Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, eris.Wrap(err, "could not dial rabbitmq")
	}
	return conn, nil
}

// CreateConsumer implements AsyncMessagingPort.
func (r *RabbitMQClient) CreateConsumer() (Consumer, error) {
	if r.config.ConsumerName == "" {
		return nil, ErrConsumerNameIsRequired
	}
	if r.config.QueueName == "" {
		return nil, ErrQueueNameIsRequired
	}
	if r.config.ExchangeName == "" {
		return nil, ErrExchangeNameIsRequired
	}

	conn, err := createConnection(r.config)
	if err != nil {
		return nil, eris.Wrap(err, "could not create connection for consumer")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, eris.Wrap(err, "could not create channel for consumer")
	}

	return &RabbitMQConsumer{
		config:     r.config,
		controller: NewChannelController(r.config.Logger, conn, ch, r.config),
		wg:         &sync.WaitGroup{},
	}, nil
}

// CreatePublisher implements AsyncMessagingPort.
func (r *RabbitMQClient) CreatePublisher() (Publisher, error) {
	if r.config.ExchangeName == "" {
		return nil, ErrExchangeNameIsRequired
	}

	conn, err := createConnection(r.config)
	if err != nil {
		return nil, eris.Wrap(err, "could not create connection for publisher")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, eris.Wrap(err, "could not create channel for publisher")
	}

	if r.config.IsConfirmMode {
		err = ch.Confirm(false)
		if err != nil {
			return nil, eris.Wrap(err, "could not set channel to confirm mode")
		}
	}

	return &RabbitMQPublisher{
		config:              r.config,
		controller:          NewChannelController(r.config.Logger, conn, ch, r.config),
		unpublishedMessages: make(chan messaging.Message),
	}, nil
}

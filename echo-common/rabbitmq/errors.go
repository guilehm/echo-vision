package rabbitmq

import "errors"

var (
	ErrCouldNotDecodeMessage  = errors.New("could not decode message")
	ErrExchangeNameIsRequired = errors.New("exchange name is required")
	ErrQueueNameIsRequired    = errors.New("queue name is required")
	ErrConsumerNameIsRequired = errors.New("consumer name is required")
	ErrCouldNotCreateChannel  = errors.New("could not create channel")
	ErrRabbitMQURLIsRequired  = errors.New("RabbitMQ URL is required")
	ErrorLoggerIsRequired     = errors.New("Logger is required")
)

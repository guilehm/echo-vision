package rabbitmq

import (
	"context"
)

// Message defines the structure of a message handled by the async messaging system.
type Message struct {
	Key     string
	Payload []byte
	Headers map[string]string
}

// AsyncMessagingPort defines the interface for asynchronous messaging.
type AsyncMessagingPort interface {
	CreatePublisher() (Publisher, error)
	CreateConsumer() (Consumer, error)
	Close() error
}

// type EventDecoder interface {
// 	Decode(v interface{}) error
// }

type Handler interface {
	Topics() []string
	Handle(ctx context.Context, topic Topic, msg Message) HandlerResponse
}

type Publisher interface {
	Publish(ctx context.Context, topic string, message Message) error
	Close() error
}

type Consumer interface {
	Subscribe(ctx context.Context, topic string, handler func(msg Message) error) error
	Close() error
}

type HandlerResponse int

const (
	Success HandlerResponse = iota
	DeadLetter
)

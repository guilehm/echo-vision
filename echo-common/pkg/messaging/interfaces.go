package messaging

import (
	"context"
)

// AsyncMessagingPort defines the interface for asynchronous messaging.
type AsyncMessagingPort interface {
	CreatePublisher() (Publisher, error)
	CreateConsumer() (Consumer, error)
}

// Publisher defines the interface for a message publisher.
type Publisher interface {
	StartPublisher(ctx context.Context) error
	Publish(ctx context.Context, msg Message) error
	Close() error
}

// Consumer defines the interface for a message consumer.
type Consumer interface {
	Subscribe(ctx context.Context, handler Handler) error
	Close() error
}

// Handler defines the interface for a message handler.
type Handler interface {
	Topics() []string
	Handle(ctx context.Context, msg Message) HandlerResponse
}

// HandlerResponse defines the possible responses from a message handler.
type HandlerResponse int

const (
	Success HandlerResponse = iota
	DeadLetter
)

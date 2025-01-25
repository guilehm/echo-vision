package messaging

import "context"

// Message defines the structure of a message handled by the async messaging system.
type Message struct {
	Topic   string
	Payload []byte
	Headers map[string]string
}

// Publisher defines the interface for a message publisher.
type Publisher interface {
	Publish(ctx context.Context, msg Message) error
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

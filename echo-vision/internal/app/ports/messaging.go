package ports

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
	CreateChannel() (MessagingChannel, error)
	Close() error
}

// MessagingChannel defines operations that can be performed on a specific channel.
type MessagingChannel interface {
	Publish(ctx context.Context, topic string, message Message) error
	Subscribe(ctx context.Context, topic string, handler func(msg Message) error) error
	Close() error
}

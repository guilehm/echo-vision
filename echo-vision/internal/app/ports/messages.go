package ports

import (
	"context"
)

type Message struct {
	Key     string
	Payload []byte
	Headers map[string]string
}

type AsyncMessagingPort interface {
	Publish(ctx context.Context, topic string, message Message) error
	Subscribe(ctx context.Context, topic string, handler func(msg Message) error) error
	Close() error
}

package rabbitmqmocks

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-common/rabbitmq"
)

// mockConsumer is a mock implementation of messaging.Consumer.
type mockConsumer struct {
	mockedChan chan messaging.Message
}

// Close implements messaging.Consumer.
func (m *mockConsumer) Close() error {
	return nil
}

// Subscribe implements messaging.Consumer.
func (m *mockConsumer) Subscribe(ctx context.Context, handler messaging.Handler) error {
	for msg := range m.mockedChan {
		go handler.Handle(ctx, msg)
	}
	return nil
}

// TODO: use interface from pkg to implement this mock
// NewConsumer creates a new mock consumer.
func NewConsumer() rabbitmq.Consumer {
	return &mockConsumer{
		mockedChan: make(chan messaging.Message),
	}
}

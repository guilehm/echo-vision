package rabbitmqmocks

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
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

// NewConsumer creates a new mock consumer.
func NewConsumer(messageChan chan messaging.Message) messaging.Consumer {
	return &mockConsumer{
		mockedChan: messageChan,
	}
}

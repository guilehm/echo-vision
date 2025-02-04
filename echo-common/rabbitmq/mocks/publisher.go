package rabbitmqmocks

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

// mockPublisher is a mock implementation of messaging.Publisher.
type mockPublisher struct {
	mockedChan chan messaging.Message
}

// Close implements messaging.Publisher.
func (m *mockPublisher) Close() error {
	return nil
}

// Publish implements messaging.Publisher.
func (m *mockPublisher) Publish(ctx context.Context, msg messaging.Message) error {
	m.mockedChan <- msg
	return nil
}

// StartPublisher implements rabbitmq.Publisher.
func (m *mockPublisher) StartPublisher(ctx context.Context) error {
	return nil
}

// NewPublisher creates a new mock publisher.
func NewPublisher(messageChan chan messaging.Message) messaging.Publisher {
	return &mockPublisher{
		mockedChan: messageChan,
	}
}

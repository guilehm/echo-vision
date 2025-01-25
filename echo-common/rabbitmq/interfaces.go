package rabbitmq

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
)

// AsyncMessagingPort defines the interface for asynchronous messaging.
type AsyncMessagingPort interface {
	CreatePublisher() (Publisher, error)
	CreateConsumer() (Consumer, error)
}


// Publisher defines the interface for a message publisher.
type Publisher interface {
	StartPublisher(ctx context.Context) error
	Publish(ctx context.Context, msg messaging.Message) error
	Close() error
}

// Consumer defines the interface for a message consumer.
type Consumer interface {
	Subscribe(ctx context.Context, handler messaging.Handler) error
	Close() error
}


package rabbitmqmocks

import (
	"context"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/stretchr/testify/mock"
)

// Handler is a mock implementation of rabbitmq.Handler.
type Handler struct {
	Mock mock.Mock
}

// Handle implements rabbitmq.Handler.
func (m *Handler) Handle(ctx context.Context, msg messaging.Message) messaging.HandlerResponse {
	ret := m.Mock.Called(ctx, msg)
	var r0 messaging.HandlerResponse
	if rf, ok := ret.Get(0).(func(context.Context, messaging.Message) messaging.HandlerResponse); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Get(0).(messaging.HandlerResponse)
	}
	return r0
}

// Topics implements rabbitmq.Handler.
func (m *Handler) Topics() []string {
	return []string{}
}

// ResetMock resets the mock.
func (m *Handler) ResetMock() {
	m.Mock = mock.Mock{}
}

// NewHandler creates a new mock handler.
func NewHandler() *Handler {
	return &Handler{
		Mock: mock.Mock{},
	}
}

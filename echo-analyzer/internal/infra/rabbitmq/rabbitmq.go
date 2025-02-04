package rabbitmqadapter

import (
	"context"
	"fmt"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"

	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type RabbitMQAdapter struct{}

func (r *RabbitMQAdapter) Topics() []string {
	return []string{
		hubevents.EventImageAnalysCreated,
		hubevents.EventImageAnalysisStatusUpdated,
	}
}

func (r *RabbitMQAdapter) Handle(ctx context.Context, msg messaging.Message) messaging.HandlerResponse {
	switch msg.Topic {
	case hubevents.EventImageAnalysCreated:
		fmt.Println("image analysis from ECHO-ANALYZER", msg)
		return messaging.Success // case hubevents.EventImageAnalysisStatusUpdated:
	// 	fmt.Println("Image analysis status updated")
	// 	return messaging.Success
	default:
		return messaging.DeadLetter
	}
}

func NewRabbitMQAdapter() *RabbitMQAdapter {
	return &RabbitMQAdapter{}
}

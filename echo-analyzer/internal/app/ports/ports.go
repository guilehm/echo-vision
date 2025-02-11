package ports

import (
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type ConsumerPort interface {
	ProcessImageAnalysis(topic string, message hubevents.EventMessage) messaging.HandlerResponse
}

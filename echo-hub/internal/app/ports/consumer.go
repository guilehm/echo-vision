package ports

import (
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type ConsumerPort interface {
	ImageAnalysisStatusUpdate(topic string, message hubevents.EventStatusUpdateMessage) messaging.HandlerResponse
}

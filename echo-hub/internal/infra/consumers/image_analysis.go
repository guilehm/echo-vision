package consumers

import (
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// ImageAnalysisStatusUpdate implements ports.ConsumerPort.
func (c *ConsumerGroup) ImageAnalysisStatusUpdate(topic string, message hubevents.EventStatusUpdateMessage) messaging.HandlerResponse {
	panic("unimplemented")
}

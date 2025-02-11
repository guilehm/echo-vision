package consumers

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// ProcessImageAnalysis implements ports.ConsumerPort.
func (c *ConsumerGroup) ProcessImageAnalysis(topic string, message hubevents.EventMessage) messaging.HandlerResponse {
	labels, err := c.irs.DetectLabels(message.File.Filepath)
	if err != nil {
		logger.Error("could not detect labels: ", slog.String("error", err.Error()))
		return messaging.DeadLetter
	}

	d, _ := json.MarshalIndent(labels, "", "  ")
	fmt.Println(string(d))

	logger.Info("labels successfully detected")
	// TODO: publish event processed message

	return messaging.Success
}

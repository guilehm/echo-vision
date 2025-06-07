package consumers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// ImageAnalysisStatusUpdate implements ports.ConsumerPort.
func (c *ConsumerGroup) ImageAnalysisStatusUpdate(topic string, id uuid.UUID, status string) messaging.HandlerResponse {
	ctx := context.Background()
	err := c.EventUseCase.SetEventStatus(ctx, id, hubevents.EventStatus(status))
	if err != nil {
		logger.Error("could not set event status",
			slog.String("id", id.String()),
			slog.String("status", status),
			slog.String("error", err.Error()),
		)
		return messaging.DeadLetter
	}

	logger.Info("event status updated",
		slog.String("id", id.String()),
		slog.String("status", status))
	return messaging.Success
}

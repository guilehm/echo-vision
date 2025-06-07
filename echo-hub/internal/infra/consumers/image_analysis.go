package consumers

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

// ImageAnalysisStatusUpdate implements ports.ConsumerPort.
func (c *ConsumerGroup) ImageAnalysisStatusUpdate(id uuid.UUID, status string, result json.RawMessage) messaging.HandlerResponse {
	ctx := context.Background()
	err := c.EventUseCase.HandleEventStatusUpdate(ctx, id, hubevents.EventStatus(status), result)
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

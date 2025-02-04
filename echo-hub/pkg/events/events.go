package hubevents

import (
	"fmt"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
)

const (
	EventImageAnalysCreated = "hub.event.image_analysis.created"
)

const (
	EventImageAnalysisStatusUpdated           = "hub.event.image_analysis.status_updated"
	EventImageAnalysisStatusUpdatedPending    = "hub.event.image_analysis.status_updated.pending"
	EventImageAnalysisStatusUpdatedProcessing = "hub.event.image_analysis.status_updated.processing"
	EventImageAnalysisStatusUpdatedCompleted  = "hub.event.image_analysis.status_updated.completed"
	EventImageAnalysisStatusUpdatedFailed     = "hub.event.image_analysis.status_updated.failed"
)

func BuildEventCreatedTopic(eventType domain.EventType) string {
	return fmt.Sprintf("hub.event.%s.created", eventType)
}

func BuildEventStatusUpdatedTopic(eventType domain.EventType, status domain.EventStatus) string {
	return fmt.Sprintf("hub.event.%s.status_updated.%s", eventType, status)
}

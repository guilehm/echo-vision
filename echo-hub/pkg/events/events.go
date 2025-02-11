package hubevents

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
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

type EventMessage struct {
	ID      uuid.UUID           `json:"id"`
	Type    domain.EventType    `json:"type"`
	SubType domain.EventSubType `json:"subType"`
	File    *valueobjects.File  `json:"file"`
}

func (e *EventMessage) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

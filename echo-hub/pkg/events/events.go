package hubevents

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
)

type EventType string

func (et EventType) Values() []EventType {
	return []EventType{
		EventTypeImageAnalysis,
	}
}

func (et EventType) StringValues() []string {
	return toStringValues(et.Values())
}

func (et EventType) String() string {
	return string(et)
}

func (et EventType) IsValid() bool {
	return isIn(et, et.Values())
}

const (
	EventTypeImageAnalysis EventType = "image_analysis"
)

type EventSubType string

func (est EventSubType) Values() []EventSubType {
	return []EventSubType{
		EventSubTypeDetectLabels,
		EventSubTypeDetectFaces,
	}
}

func (est EventSubType) StringValues() []string {
	return toStringValues(est.Values())
}

func (est EventSubType) String() string {
	return string(est)
}

const (
	EventSubTypeDetectLabels EventSubType = "detect_labels"
	EventSubTypeDetectFaces  EventSubType = "detect_faces"
)

type EventStatus string

const (
	EventStatusPending    EventStatus = "pending"
	EventStatusProcessing EventStatus = "processing"
	EventStatusCompleted  EventStatus = "completed"
	EventStatusFailed     EventStatus = "failed"
)

func (es EventStatus) String() string {
	return string(es)
}

func (es EventStatus) Values() []EventStatus {
	return []EventStatus{
		EventStatusPending,
		EventStatusProcessing,
		EventStatusCompleted,
		EventStatusFailed,
	}
}

func (es EventStatus) StringValues() []string {
	return toStringValues(es.Values())
}

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

func BuildEventCreatedTopic(eventType EventType) string {
	return fmt.Sprintf("hub.event.%s.created", eventType)
}

func BuildEventStatusUpdatedTopic(eventType EventType, status EventStatus) string {
	return fmt.Sprintf("hub.event.%s.status_updated.%s", eventType, status)
}

type EventMessage struct {
	ID      uuid.UUID          `json:"id"`
	Type    EventType          `json:"type"`
	SubType EventSubType       `json:"subType"`
	File    *valueobjects.File `json:"file"`
}

func (e *EventMessage) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

type EventStatusUpdateMessage struct {
	ID     uuid.UUID       `json:"id"`
	Type   EventType       `json:"type"`
	Status EventStatus     `json:"status"`
	Data   json.RawMessage `json:"data,omitempty"`
}

func (e *EventStatusUpdateMessage) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func toStringValues[T ~string](values []T) []string {
	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = string(v)
	}
	return stringValues
}

func isIn[T comparable](s T, values []T) bool {
	return slices.Contains(values, s)
}

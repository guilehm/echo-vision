package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

type EventType string

const (
	EventTypeImageAnalysis EventType = "image_analysis"
)

func (et EventType) Values() []EventType {
	return []EventType{
		EventTypeImageAnalysis,
	}
}

func (et EventType) StringValues() []string {
	return toStringValues(et.Values())
}

type EventSubType string

func (est EventSubType) Values() []EventSubType {
	return []EventSubType{
		EventSubTypeDetecLabels,
	}
}

func (est EventSubType) StringValues() []string {
	return toStringValues(est.Values())
}

const (
	EventSubTypeDetecLabels EventSubType = "detect_labels"
)

type EventStatus string

const (
	EventStatusPending    EventStatus = "pending"
	EventStatusProcessing EventStatus = "processing"
	EventStatusCompleted  EventStatus = "completed"
	EventStatusFailed     EventStatus = "failed"
)

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

type Event struct {
	id        uuid.UUID
	eventType EventType
	subType   EventSubType
	status    EventStatus
	payload   json.RawMessage
	result    json.RawMessage
	createdAt time.Time
	updatedAt time.Time
}

func NewEvent(
	id uuid.UUID,
	eventType EventType,
	subType EventSubType,
	payload json.RawMessage,
	result json.RawMessage,
) *Event {
	return &Event{
		id:        id,
		eventType: eventType,
		subType:   subType,
		status:    EventStatusPending,
		payload:   payload,
		result:    result,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (e *Event) Validate() error {
	if e.id == uuid.Nil {
		return shared.ErrInvalidID
	}
	if e.status == "" {
		return shared.ErrInvalidStatus
	}
	if e.eventType == "" {
		return shared.ErrInvalidEventType
	}
	if e.payload == nil {
		return shared.ErrInvalidPayload
	}
	return nil
}

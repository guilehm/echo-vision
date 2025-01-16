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
		EventSubTypeDetectLabels,
	}
}

func (est EventSubType) StringValues() []string {
	return toStringValues(est.Values())
}

const (
	EventSubTypeDetectLabels EventSubType = "detect_labels"
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
	user      *User
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
	user *User,
	id uuid.UUID,
	eventType EventType,
	subType EventSubType,
	payload json.RawMessage,
	result json.RawMessage,
	status EventStatus,
	createdAt, updatedAt time.Time,
) *Event {
	return &Event{
		user:      user,
		id:        id,
		eventType: eventType,
		subType:   subType,
		status:    status,
		payload:   payload,
		result:    result,
		createdAt: createdAt,
		updatedAt: updatedAt,
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

func (e *Event) ID() uuid.UUID {
	return e.id
}

func (e *Event) User() *User {
	return e.user
}

func (e *Event) EventType() EventType {
	return e.eventType
}

func (e *Event) SubType() EventSubType {
	return e.subType
}

func (e *Event) Status() EventStatus {
	return e.status
}

func (e *Event) Payload() json.RawMessage {
	return e.payload
}

func (e *Event) Result() json.RawMessage {
	return e.result
}

func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) UpdatedAt() time.Time {
	return e.updatedAt
}

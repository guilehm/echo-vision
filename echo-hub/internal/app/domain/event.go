package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
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

func (et EventType) String() string {
	return string(et)
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

func (est EventSubType) String() string {
	return string(est)
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

type Event struct {
	userID    uuid.UUID
	id        uuid.UUID
	eventType EventType
	subType   EventSubType
	status    EventStatus
	result    json.RawMessage
	createdAt time.Time
	updatedAt time.Time
	file      *valueobjects.File
}

func NewEvent(
	userID uuid.UUID,
	id uuid.UUID,
	eventType EventType,
	subType EventSubType,
	result json.RawMessage,
	status EventStatus,
	file *valueobjects.File,
	createdAt, updatedAt time.Time,
) *Event {
	return &Event{
		userID:    userID,
		id:        id,
		eventType: eventType,
		subType:   subType,
		status:    status,
		result:    result,
		file:      file,
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
	if err := e.validateTypes(); err != nil {
		return err
	}
	return nil
}

func (e *Event) validateTypes() error {
	if !isIn(e.eventType, e.eventType.Values()) {
		return shared.ErrInvalidEventType
	}
	if !isIn(e.subType, e.subType.Values()) {
		return shared.ErrInvalidSubType
	}
	if !isIn(e.status, e.status.Values()) {
		return shared.ErrInvalidStatus
	}
	return nil
}

func (e *Event) ID() uuid.UUID {
	return e.id
}

func (e *Event) UserID() uuid.UUID {
	return e.userID
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

func (e *Event) Result() json.RawMessage {
	return e.result
}

func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Event) File() *valueobjects.File {
	return e.file
}

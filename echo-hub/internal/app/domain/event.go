package domain

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type Event struct {
	userID    uuid.UUID
	id        uuid.UUID
	eventType hubevents.EventType
	subType   hubevents.EventSubType
	status    hubevents.EventStatus
	result    json.RawMessage
	createdAt time.Time
	updatedAt time.Time
	file      *valueobjects.File
}

func NewEvent(
	userID uuid.UUID,
	id uuid.UUID,
	eventType hubevents.EventType,
	subType hubevents.EventSubType,
	result json.RawMessage,
	status hubevents.EventStatus,
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

func (e *Event) EventType() hubevents.EventType {
	return e.eventType
}

func (e *Event) SubType() hubevents.EventSubType {
	return e.subType
}

func (e *Event) Status() hubevents.EventStatus {
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

func isIn[T comparable](s T, values []T) bool {
	return slices.Contains(values, s)
}

func (e *Event) SetStatus(status hubevents.EventStatus) error {
	if !isIn(status, e.status.Values()) {
		return shared.ErrInvalidStatus
	}
	// TODO: add state validation here
	e.status = status
	e.updatedAt = time.Now()
	return nil
}

func (e *Event) SetResult(result json.RawMessage) {
	e.result = result
	e.updatedAt = time.Now()
}

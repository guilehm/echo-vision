package tests

import (
	"context"
	"encoding/json"
	"time"

	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
)

// TODO: remove validUser and validEvent
var validUser = domain.NewUser(
	uuid.New(),
	"Arthur",
	"Morgan",
	"arthur@gmail.com",
	time.Now(),
	time.Now(),
)

var validEvent = domain.NewEvent(
	validUser.ID(),
	uuid.New(),
	hubevents.EventTypeImageAnalysis,
	hubevents.EventSubTypeDetectLabels,
	json.RawMessage(`{"result": "success"}`),
	hubevents.EventStatusPending,
	nil,
	time.Now(),
	time.Now(),
)

var validEventWithFile = domain.NewEvent(
	validUser.ID(),
	uuid.New(),
	hubevents.EventTypeImageAnalysis,
	hubevents.EventSubTypeDetectLabels,
	json.RawMessage(`{"result": "success"}`),
	hubevents.EventStatusPending,
	valueobjects.NewFile(
		"path/to/file.jpg",
		"file.jpg",
		"image/jpeg",
		1024,
	),
	time.Now(),
	time.Now(),
)

var ctx = context.Background()

func makeUser(email string) *domain.User {
	return domain.NewUser(
		uuid.New(),
		"Arthur",
		"Morgan",
		email,
		time.Now(),
		time.Now(),
	)
}

func makeEvent(u *domain.User) *domain.Event {
	return domain.NewEvent(
		u.ID(),
		uuid.New(),
		hubevents.EventTypeImageAnalysis,
		hubevents.EventSubTypeDetectLabels,
		json.RawMessage(`{"result": "success"}`),
		hubevents.EventStatusPending,
		nil,
		time.Now(),
		time.Now(),
	)
}

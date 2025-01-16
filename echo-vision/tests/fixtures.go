package tests

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
)

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
	domain.EventTypeImageAnalysis,
	domain.EventSubTypeDetectLabels,
	json.RawMessage(`{"key": "value"}`),
	json.RawMessage(`{"result": "success"}`),
	domain.EventStatusPending,
	time.Now(),
	time.Now(),
)

var ctx = context.Background()

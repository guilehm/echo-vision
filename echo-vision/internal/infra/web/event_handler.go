package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
)

type EventHandler struct {
	eventPort ports.EventPort
}

func NewEventHandler(up ports.EventPort) *EventHandler {
	return &EventHandler{
		eventPort: up,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input ports.EventCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("error decoding request body", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			"error decoding request body",
		)))
		return
	}
	d, err := json.Marshal(input)
	if err != nil {
		logger.Error("error marshalling input", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			"error marshalling input",
		)))
		return
	}

	now := time.Now()
	event := domain.NewEvent(
		uuid.New(), // TODO: get user id from context
		uuid.New(),
		domain.EventType(input.EventType),
		domain.EventSubType(input.SubType),
		d,
		nil,
		domain.EventStatusPending,
		now,
		now,
	)
	if err := event.Validate(); err != nil {
		logger.Error("error validating event", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			err.Error(),
		)))
		return
	}

	eventResp, err := h.eventPort.SaveEvent(r.Context(), event)
	if err != nil {
		logger.Error("error creating event", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse[uuid.UUID](&eventResp, nil))
}

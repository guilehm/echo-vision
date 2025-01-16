package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

type EventHandler struct {
	eventPort ports.EventPort
}

// ListEvents implements ports.EventWebPort.
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := chi.URLParam(r, "userID")
	user, err := fromContext[domain.User](ctx, contextKeyMeUser)
	if err != nil {
		logger.Error("error getting user from context", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusForbidden,
			err.Error(),
		)))
		return
	}

	// TODO: we may change this behaviour in the future
	if userID != user.ID().String() {
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)))
		return
	}

	events, err := h.eventPort.EventsByUser(ctx, user.ID())
	if err != nil {
		logger.Error("error getting events by user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse[ports.ApiListResponse[ports.EventResponse]](&ports.ApiListResponse[ports.EventResponse]{
		Results: ports.MapEventsToApiResponse(events),
	}, nil))
}

func NewEventHandler(up ports.EventPort) ports.EventWebPort {
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
			shared.ErrDecodingRequestBody.Error(),
		)))
		return
	}

	ctx := r.Context()

	user, err := fromContext[domain.User](ctx, contextKeyMeUser)
	if err != nil {
		logger.Error("error getting user from context", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusForbidden,
			err.Error(),
		)))
		return
	}

	event, err := h.eventPort.CreateEvent(ctx, user.ID(), input.EventType, input.SubType)
	if err != nil {
		logger.Error("error creating event", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			err.Error(),
		)))
		return
	}

	eventID, err := h.eventPort.SaveEvent(ctx, event)
	if err != nil {
		logger.Error("error saving event", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse[ports.EventCreateResponse](&ports.EventCreateResponse{
		ID: eventID,
	}, nil))
}

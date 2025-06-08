package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain/valueobjects"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
)

// EventHandler is a handler for events.
type EventHandler struct {
	eventPort ports.EventPort
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler(up ports.EventPort) ports.EventWebPort {
	return &EventHandler{
		eventPort: up,
	}
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

	cursor := r.URL.Query().Get("cursor")
	limitQuery := r.URL.Query().Get("limit")
	limit := 10
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil || limit <= 0 || limit > 100 {
			logger.Error("invalid limit parameter", slog.String("limit", limitQuery))
			handleApiResponse(w, apiResponse[any](nil, newApiError(
				http.StatusBadRequest,
				"invalid limit parameter",
			)))
			return
		}
	}

	events, nextCursor, err := h.eventPort.EventsByUser(ctx, user.ID(), limit, cursor)
	if err != nil {
		logger.Error("error getting events by user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse(&ports.ApiListResponse[dtos.EventResponse]{
		Results:    ports.MapEventsToApiResponse(events),
		NextCursor: nextCursor,
	}, nil))
}

// ListOwnEvents implements ports.EventWebPort.
func (h *EventHandler) ListOwnEvents(w http.ResponseWriter, r *http.Request) {
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

	cursor := r.URL.Query().Get("cursor")
	limitQuery := r.URL.Query().Get("limit")
	limit := 10
	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil || limit <= 0 || limit > 100 {
			logger.Error("invalid limit parameter", slog.String("limit", limitQuery))
			handleApiResponse(w, apiResponse[any](nil, newApiError(
				http.StatusBadRequest,
				"invalid limit parameter",
			)))
			return
		}
	}
	events, nextCursor, err := h.eventPort.EventsByUser(ctx, user.ID(), limit, cursor)
	if err != nil {
		logger.Error("error getting events by user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse(&ports.ApiListResponse[dtos.EventResponse]{
		Results:    ports.MapEventsToApiResponse(events),
		NextCursor: nextCursor,
	}, nil))
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

	var file *valueobjects.File
	if input.Filepath != "" {
		file = valueobjects.NewFile(
			input.Filepath,
			input.Filename,
			input.ContentType,
			input.Filesize,
		)
		if !file.IsValid() {
			logger.Error("error creating file", slog.String("error", shared.ErrInvalidFile.Error()))
			handleApiResponse(w, apiResponse[any](nil, newApiError(
				http.StatusBadRequest,
				shared.ErrInvalidFile.Error(),
			)))
			return
		}
	}

	event, err := h.eventPort.CreateEvent(
		ctx,
		user.ID(),
		input.EventType,
		input.SubType,
		file,
	)
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

	handleApiResponse(w, apiResponse(&dtos.EventCreateResponse{
		ID: eventID,
	}, nil))
}

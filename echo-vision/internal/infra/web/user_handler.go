package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
)

type UserHandler struct {
	userPort ports.UserPort
}

func NewUserHandler(up ports.UserPort) *UserHandler {
	return &UserHandler{
		userPort: up,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input ports.UserCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("error decoding request body", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			"error decoding request body",
		)))
		return
	}

	var err error
	var user *domain.User

	if user, err = h.userPort.CreateUser(
		r.Context(),
		input.FirstName,
		input.LastName,
		input.Email,
		input.Password,
	); err != nil {
		logger.Error("error creating user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			err.Error(),
		)))
		return

	}

	userID, err := h.userPort.SaveUser(r.Context(), user)
	if err != nil {
		logger.Error("error saving user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse[ports.UserCreateResponse](&ports.UserCreateResponse{
		ID:           userID,
		AccessToken:  user.AccessToken(),
		RefreshToken: user.RefreshToken(),
	}, nil))
}

func (h *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	panic("FindUserByID not implemented")
}

func (h *UserHandler) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	panic("FindUserByEmail not implemented")
}

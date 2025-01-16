package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/shared"
)

type UserHandler struct {
	userPort ports.UserPort
}

func NewUserHandler(up ports.UserPort) ports.UserWebPort {
	return &UserHandler{
		userPort: up,
	}
}

// CreateUser implements ports.UserWebPort.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input ports.UserCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("error decoding request body", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			shared.ErrDecodingRequestBody.Error(),
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

// Login implements ports.UserWebPort.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input ports.UserLoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("error decoding request body", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			shared.ErrDecodingRequestBody.Error(),
		)))
		return
	}
	user, err := h.userPort.AuthenticateUser(r.Context(), input.Email, input.Password)
	if err != nil {
		logger.Error("error authenticating user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse[ports.UserLoginResponse](&ports.UserLoginResponse{
		ID:           user.ID(),
		AccessToken:  user.AccessToken(),
		RefreshToken: user.RefreshToken(),
	}, nil))
}

// MeUser implements ports.UserWebPort.
func (h *UserHandler) MeUser(w http.ResponseWriter, r *http.Request) {
	u, err := fromContext[domain.User](r.Context(), contextKeyMeUser)
	if err != nil {
		logger.Error("error getting user from context", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusUnauthorized,
			"unauthorized",
		)))
		return
	}
	handleApiResponse(w, apiResponse[ports.UserResponse](&ports.UserResponse{
		ID:        u.ID(),
		Email:     u.Email(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
	}, nil))
}

// Logout implements ports.UserWebPort.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// RefreshToken implements ports.UserWebPort.
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports/dtos"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
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

	_, err = h.userPort.SaveUser(r.Context(), user)
	if err != nil {
		logger.Error("error saving user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	// TODO: set cookies

	handleApiResponse(w, apiResponse(&dtos.UserCreateResponse{
		ID:           user.ID(),
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

	// TODO: setup cookie expiration
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    user.AccessToken(),
		Path:     "/",
		HttpOnly: false,
		// TODO: control based on env
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "refreshToken",
		Value: user.RefreshToken(),
		// TODO: do not use / as path for refresh token
		Path:     "/",
		HttpOnly: false,
		// TODO: control based on env
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	handleApiResponse(w, apiResponse(&dtos.UserLoginResponse{
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
	handleApiResponse(w, apiResponse[dtos.UserResponse](ports.MapUserToApiResponse(u), nil))
}

// Logout implements ports.UserWebPort.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		// TODO: control based on env
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "refreshToken",
		Value: "",
		// TODO: do not use / as path for refresh token
		Path:     "/",
		HttpOnly: false,
		// TODO: control based on env
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	handleApiResponse(w, apiResponse[any](nil, nil))
}

// RefreshToken implements ports.UserWebPort.
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input ports.UserRefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("error decoding request body", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			shared.ErrDecodingRequestBody.Error(),
		)))
		return
	}
	if input.RefreshToken == "" {
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)))
		return
	}
	u, err := h.userPort.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		logger.Error("error refreshing token", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}
	handleApiResponse(w, apiResponse[dtos.RefreshTokenResponse](&dtos.RefreshTokenResponse{
		AccessToken: u.AccessToken(),
	}, nil))
}

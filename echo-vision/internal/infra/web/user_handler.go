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
		handleApiResponse(w, apiResponse[any](r.Context(), nil, apiError{
			Status:  http.StatusBadRequest,
			Message: "error decoding request body",
		}))
		return
	}

	now := time.Now()
	user := domain.NewUser(
		uuid.New(),
		input.FirstName,
		input.LastName,
		input.Email,
		now,
		now,
	)
	if err := user.Validate(); err != nil {
		logger.Error("error validating user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](r.Context(), nil, apiError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}))
		return
	}

	userResp, err := h.userPort.SaveUser(r.Context(), user)
	if err != nil {
		logger.Error("error saving user", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](r.Context(), nil, err))
		return
	}
	handleApiResponse(w, apiResponse[uuid.UUID](r.Context(), &userResp, nil))
}

func (h *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	panic("FindUserByID not implemented")
	// id := chi.URLParam(r, "id")
	//
	// userID, err := uuid.Parse(id)
	// if err != nil {
	// 	logger.Info("error parsing user id", slog.String("error", err.Error()))
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	//
	// userResp, err := h.userPort.FindUserByID(r.Context(), userID)
	// if err != nil {
	// 	if errors.Is(err, shared.ErrNotNound) {
	// 		http.Error(w, err.Error(), http.StatusNotFound)
	// 		return
	// 	}
	// 	logger.Error("error finding user by id", slog.String("error", err.Error()))
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// writeJson(w, userResp)
}

func (h *UserHandler) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	panic("FindUserByEmail not implemented")
	// email := chi.URLParam(r, "email")
	//
	// userResp, err := h.userPort.FindUserByEmail(r.Context(), email)
	// if err != nil {
	// 	if errors.Is(err, shared.ErrNotNound) {
	// 		http.Error(w, err.Error(), http.StatusNotFound)
	// 		return
	// 	}
	// 	logger.Error("error finding user by email", slog.String("error", err.Error()))
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	//
	// writeJson(w, userResp)
}

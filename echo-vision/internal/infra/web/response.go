package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guilehm/echo-vision/internal/app/shared"
	"github.com/lib/pq"
)

type ApiResponse[T any] struct {
	Error  string `json:"error,omitempty"`
	Data   *T     `json:"data,omitempty"`
	Status int    `json:"status,omitempty"`
}

func apiResponse[T any](data *T, err error) *ApiResponse[T] {
	var errorMessage string
	var status int
	if err != nil {
		// set default error status and message
		errorMessage = err.Error()
		status = http.StatusInternalServerError

		// handle not found errors
		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
			errorMessage = http.StatusText(http.StatusNotFound)
		}

		// user not found returns bad request for security
		if errors.Is(err, shared.ErrUserNotFound) {
			status = http.StatusBadRequest
			errorMessage = http.StatusText(http.StatusNotFound)
		}

		if errors.Is(err, shared.ErrInvalidPassword) {
			status = http.StatusBadRequest
			errorMessage = http.StatusText(http.StatusBadRequest)
		}

		// handle specific error type for postgres
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			status, errorMessage = logPGError(pgErr)
		}

		// handle specific api error type
		var apiErr *apiError
		if errors.As(err, &apiErr) {
			status = apiErr.Status
			errorMessage = apiErr.Message
		}

		return &ApiResponse[T]{
			Data:   data,
			Error:  errorMessage,
			Status: status,
		}
	}

	return &ApiResponse[T]{
		Data:   data,
		Status: http.StatusOK,
	}
}

func handleApiResponse[T any](w http.ResponseWriter, response *ApiResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	_ = json.NewEncoder(w).Encode(response)
}

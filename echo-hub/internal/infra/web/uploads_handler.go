package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
)

type UploadHandler struct {
	filestoragePort filestorage.FileStoragePort
}

func NewUploadHandler(uph filestorage.FileStoragePort) *UploadHandler {
	return &UploadHandler{
		filestoragePort: uph,
	}
}

func (h *UploadHandler) PresignedURL(w http.ResponseWriter, r *http.Request) {
	var input ports.UploadPresignedURLInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			shared.ErrDecodingRequestBody.Error(),
		)))
		return
	}

	if !input.IsValid() {
		handleApiResponse(w, apiResponse[any](nil, newApiError(
			http.StatusBadRequest,
			shared.ErrInvalidPayload.Error(),
		)))
		return
	}

	// TODO: validate file data

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
	url, err := h.filestoragePort.GeneratePreSignedURL(
		filestorage.NewFileKey(
			"users/"+user.ID().String(),
			input.Filepath,
		),
		input.ContentType,
	)
	if err != nil {
		logger.Error("error generating presigned url", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse(&ports.UploadPresignedURLResponse{
		URL: url,
	}, nil))
}

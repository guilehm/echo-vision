package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type UploadHandler struct {
	filestoragePort filestorage.FileStoragePort
}

func NewUploadHandler(fsp filestorage.FileStoragePort) *UploadHandler {
	return &UploadHandler{
		filestoragePort: fsp,
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

	// TODO: validate content type

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

	path := fmt.Sprintf("users/%s", user.ID().String())

	// TODO: move to a function
	switch input.EventType {
	case hubevents.EventTypeImageAnalysis.String():
		path = fmt.Sprintf("%s/%s", path, "image-analysis")
	}

	fk := filestorage.NewUploadFileKey(
		path,
		input.Filename,
	)

	url, err := h.filestoragePort.GeneratePreSignedURL(
		fk,
		input.ContentType,
	)
	if err != nil {
		logger.Error("error generating presigned url", slog.String("error", err.Error()))
		handleApiResponse(w, apiResponse[any](nil, err))
		return
	}

	handleApiResponse(w, apiResponse(&ports.UploadPresignedURLResponse{
		URL:      url,
		Filepath: fk.Filepath,
		Filename: fk.Filename,
	}, nil))
}

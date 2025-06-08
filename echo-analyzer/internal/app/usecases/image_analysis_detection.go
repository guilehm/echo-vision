package usecases

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
)

// DetectFaces implements ports.ImageAnalysisPort.
func (i ImageAnalysisUseCase) DetectFaces(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.FaceDetail, error) {
	detecFunc := func(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.FaceDetail, error) {
		faces, err := i.irs.DetectFaces(filePath)
		if err != nil {
			return nil, err
		}
		return faces, nil
	}

	faces, err := runDetection(ctx, i, detecFunc, eventID, filePath)
	if err != nil {
		logger.Error("failed to detect faces", slog.String("error", err.Error()))
		return nil, err
	}
	return faces, nil
}

// DetectLabels implements ports.ImageAnalysisPort.
func (i ImageAnalysisUseCase) DetectLabels(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.Label, error) {
	detecFunc := func(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.Label, error) {
		labels, err := i.irs.DetectLabels(filePath)
		if err != nil {
			return nil, err
		}
		return labels, nil
	}

	labels, err := runDetection(ctx, i, detecFunc, eventID, filePath)
	if err != nil {
		logger.Error("failed to detect labels", slog.String("error", err.Error()))
		return nil, err
	}

	return labels, nil
}

func runDetection[T any](
	ctx context.Context,
	i ImageAnalysisUseCase,
	detectFunc func(ctx context.Context, eventID uuid.UUID, filePath string) (T, error),
	eventID uuid.UUID,
	filePath string,
) (T, error) {
	var zeroValue T
	var err error

	defer func() {
		// publish failure status if an error occurs
		if err != nil {
			if publishErr := i.publisher.PublishImageAnalysisStatusUpdate(
				ctx,
				eventID,
				analyzerevents.EventStatusFailed,
				nil,
			); publishErr != nil {
				logger.Error("failed to publish image analysis failure status", slog.String("error", publishErr.Error()))
			}
		}
	}()

	err = i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusProcessing,
		nil,
	)
	if err != nil {
		return zeroValue, err
	}

	data, err := detectFunc(ctx, eventID, filePath)
	if err != nil {
		return zeroValue, err
	}

	err = i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusCompleted,
		toRawMessage(data),
	)
	if err != nil {
		return zeroValue, err
	}
	return data, nil
}

func toRawMessage(data any) json.RawMessage {
	raw, err := json.Marshal(data)
	if err != nil {
		logger.Error("could not marshal data to raw message", slog.String("error", err.Error()))
		return nil
	}
	return raw
}

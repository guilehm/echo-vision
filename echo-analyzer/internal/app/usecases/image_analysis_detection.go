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
	err := i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusProcessing,
		nil,
	)
	faces, err := i.irs.DetectFaces(filePath)
	if err != nil {
		return nil, err
	}
	err = i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusCompleted,
		toRawMessage(faces),
	)
	return faces, nil
}

// DetectLabels implements ports.ImageAnalysisPort.
func (i ImageAnalysisUseCase) DetectLabels(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.Label, error) {
	err := i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusProcessing,
		nil,
	)
	labels, err := i.irs.DetectLabels(filePath)
	if err != nil {
		return nil, err
	}
	err = i.publisher.PublishImageAnalysisStatusUpdate(
		ctx,
		eventID,
		analyzerevents.EventStatusCompleted,
		toRawMessage(labels),
	)
	return labels, nil
}

func toRawMessage(data any) json.RawMessage {
	raw, err := json.Marshal(data)
	if err != nil {
		logger.Error("could not marshal data to raw message", slog.String("error", err.Error()))
		return nil
	}
	return raw
}

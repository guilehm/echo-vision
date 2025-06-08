package ports

import (
	"context"

	"github.com/google/uuid"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
)

type ImageAnalysisPort interface {
	DetectLabels(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.Label, error)
	DetectFaces(ctx context.Context, eventID uuid.UUID, filePath string) ([]analyzerevents.FaceDetail, error)
}

package ports

import (
	analysistypes "github.com/guilehm/echo-vision/echo-analyzer/pkg/types"
)

// ImageRecognitionServicePort defines the interface for image recognition services.
type ImageRecognitionServicePort interface {
	BucketName() string
	DetectLabels(filepath string) ([]analysistypes.Label, error)
	DetectFaces(filepath string) ([]analysistypes.FaceDetail, error)
}

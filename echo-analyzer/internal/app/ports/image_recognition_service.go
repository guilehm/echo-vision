package ports

import (
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
)

// ImageRecognitionServicePort defines the interface for image recognition services.
type ImageRecognitionServicePort interface {
	BucketName() string
	DetectLabels(filepath string) ([]analyzerevents.Label, error)
	DetectFaces(filepath string) ([]analyzerevents.FaceDetail, error)
}

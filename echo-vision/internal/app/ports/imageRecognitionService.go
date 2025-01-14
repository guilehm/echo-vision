package ports

import "github.com/guilehm/echo-vision/internal/app/domain"

// ImageRecognitionServicePort defines the interface for image recognition services.
type ImageRecognitionServicePort interface {
	DetectLabels(imageBytes []byte) ([]domain.Label, error)
	DetectFaces(imageBytes []byte) ([]domain.FaceDetail, error)
}

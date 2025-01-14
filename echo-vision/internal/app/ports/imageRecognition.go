package ports

import "github.com/guilehm/echo-vision/internal/app/domain"

// ImageRecognitionPort defines the behavior required for image recognition.
type ImageRecognitionPort interface {
	DetectLabels(imageBytes []byte) ([]domain.Label, error)
	DetectFaces(imageBytes []byte) ([]domain.FaceDetail, error)
}

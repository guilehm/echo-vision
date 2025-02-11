package ports

import "github.com/guilehm/echo-vision/echo-analyzer/internal/app/domain"

// ImageRecognitionServicePort defines the interface for image recognition services.
type ImageRecognitionServicePort interface {
	BucketName() string
	DetectLabels(filepath string) ([]domain.Label, error)
	DetectFaces(imageInput any) ([]domain.FaceDetail, error)
}

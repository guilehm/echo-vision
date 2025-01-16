package awsrekognition

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

type AWSRekognitionAdapter struct {
	client *rekognition.Rekognition
}

// NewAWSRekognitionAdapter creates a new AWSRekognitionAdapter.
func NewAWSRekognitionAdapter(region string) (ports.ImageRecognitionServicePort, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}
	return &AWSRekognitionAdapter{
		client: rekognition.New(sess),
	}, nil
}

// DetectLabels detects labels in the provided image bytes.
func (a *AWSRekognitionAdapter) DetectLabels(imageBytes []byte) ([]domain.Label, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: imageBytes,
		},
		// TODO: make these configurable
		MaxLabels:     aws.Int64(10),
		MinConfidence: aws.Float64(70.0),
	}

	result, err := a.client.DetectLabels(input)
	if err != nil {
		return nil, err
	}

	labels := make([]domain.Label, len(result.Labels))
	for i, lbl := range result.Labels {
		labels[i] = domain.Label{
			Name:       aws.StringValue(lbl.Name),
			Confidence: aws.Float64Value(lbl.Confidence),
		}
	}
	return labels, nil
}

// DetectFaces detects faces in the provided image bytes.
func (a *AWSRekognitionAdapter) DetectFaces(imageBytes []byte) ([]domain.FaceDetail, error) {
	input := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: imageBytes,
		},
		Attributes: aws.StringSlice([]string{"ALL"}),
	}

	result, err := a.client.DetectFaces(input)
	if err != nil {
		return nil, err
	}

	faces := make([]domain.FaceDetail, len(result.FaceDetails))
	for i, face := range result.FaceDetails {
		faces[i] = domain.FaceDetail{
			BoundingBox: domain.BoundingBox{
				Top:    aws.Float64Value(face.BoundingBox.Top),
				Left:   aws.Float64Value(face.BoundingBox.Left),
				Width:  aws.Float64Value(face.BoundingBox.Width),
				Height: aws.Float64Value(face.BoundingBox.Height),
			},
			Confidence: aws.Float64Value(face.Confidence),
		}
	}
	return faces, nil
}

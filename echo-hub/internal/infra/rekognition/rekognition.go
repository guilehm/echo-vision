package awsrekognition

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/domain"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/rotisserie/eris"
)

type AWSRekognitionAdapter struct {
	client *rekognition.Client
}

// NewAWSRekognitionAdapter initializes a new AWSRekognitionAdapter.
func NewAWSRekognitionAdapter(region string) (ports.ImageRecognitionServicePort, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, eris.Wrap(err, "error loading AWS config")
	}

	return &AWSRekognitionAdapter{
		client: rekognition.NewFromConfig(cfg),
	}, nil
}

// DetectLabels detects labels in an image using AWS Rekognition.
func (a *AWSRekognitionAdapter) DetectLabels(imageBytes []byte) ([]domain.Label, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(70.0),
	}

	result, err := a.client.DetectLabels(context.TODO(), input)
	if err != nil {
		return nil, eris.Wrap(err, "error detecting labels")
	}

	labels := make([]domain.Label, len(result.Labels))
	for i, lbl := range result.Labels {
		labels[i] = domain.Label{
			Name:       lbl.Name,
			Confidence: lbl.Confidence,
		}
	}
	return labels, nil
}

// DetectFaces detects faces in an image using AWS Rekognition.
func (a *AWSRekognitionAdapter) DetectFaces(imageBytes []byte) ([]domain.FaceDetail, error) {
	input := &rekognition.DetectFacesInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
		Attributes: []types.Attribute{types.AttributeAll},
	}

	result, err := a.client.DetectFaces(context.TODO(), input)
	if err != nil {
		return nil, eris.Wrap(err, "error detecting faces")
	}

	faces := make([]domain.FaceDetail, len(result.FaceDetails))
	for i, face := range result.FaceDetails {
		faces[i] = domain.FaceDetail{
			BoundingBox: domain.BoundingBox{
				Top:    face.BoundingBox.Top,
				Left:   face.BoundingBox.Left,
				Width:  face.BoundingBox.Width,
				Height: face.BoundingBox.Height,
			},
			Confidence: face.Confidence,
		}
	}
	return faces, nil
}

package awsrekognition

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/guilehm/echo-vision/echo-analyzer/internal/app/ports"
	analyzerevents "github.com/guilehm/echo-vision/echo-analyzer/pkg/events"
	"github.com/guilehm/echo-vision/echo-common/pkg/logging"
	"github.com/rotisserie/eris"
)

var logger = logging.NewLogger()

type AWSRekognitionAdapter struct {
	client *rekognition.Client
	bucket string
}

// NewAWSRekognitionAdapter initializes a new AWSRekognitionAdapter.
func NewAWSRekognitionAdapter(region, bucketName string) (ports.ImageRecognitionServicePort, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, eris.Wrap(err, "error loading AWS config")
	}

	return &AWSRekognitionAdapter{
		client: rekognition.NewFromConfig(cfg),
		bucket: bucketName,
	}, nil
}

// DetectLabels detects labels in an image using AWS Rekognition.
func (a *AWSRekognitionAdapter) DetectLabels(filepath string) ([]analyzerevents.Label, error) {
	logger.Info("detecting labels in image", slog.String("filepath", filepath))
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(a.BucketName()),
				Name:   aws.String(filepath),
			},
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(70.0),
	}
	result, err := a.client.DetectLabels(context.TODO(), input)
	if err != nil {
		return nil, eris.Wrap(err, "error detecting labels")
	}

	labels := make([]analyzerevents.Label, len(result.Labels))
	for i, lbl := range result.Labels {
		labels[i] = labelToDomain(lbl)
	}
	logger.Info("successfully detected labels", slog.Int("count", len(labels)))
	return labels, nil
}

// DetectFaces detects faces in an image using AWS Rekognition and returns face details including emotions.
func (a *AWSRekognitionAdapter) DetectFaces(filepath string) ([]analyzerevents.FaceDetail, error) {
	logger.Info("detecting faces in image", slog.String("filepath", filepath))
	input := &rekognition.DetectFacesInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(a.BucketName()),
				Name:   aws.String(filepath),
			},
		},
		Attributes: []types.Attribute{types.AttributeAll},
	}

	result, err := a.client.DetectFaces(context.TODO(), input)
	if err != nil {
		return nil, eris.Wrap(err, "error detecting faces")
	}

	faces := make([]analyzerevents.FaceDetail, len(result.FaceDetails))
	for i, face := range result.FaceDetails {
		emotions := make([]analyzerevents.Emotion, len(face.Emotions))
		for j, emotion := range face.Emotions {
			emotions[j] = analyzerevents.Emotion{
				Type:       string(emotion.Type),
				Confidence: emotion.Confidence,
			}
		}

		faces[i] = analyzerevents.FaceDetail{
			BoundingBox: analyzerevents.BoundingBox{
				Top:    face.BoundingBox.Top,
				Left:   face.BoundingBox.Left,
				Width:  face.BoundingBox.Width,
				Height: face.BoundingBox.Height,
			},
			Confidence: face.Confidence,
			Emotions:   emotions,
		}
	}
	logger.Info("successfully detected faces", slog.Int("count", len(faces)))
	return faces, nil
}

// BucketName implements ports.ImageRecognitionServicePort.
func (a *AWSRekognitionAdapter) BucketName() string {
	return a.bucket
}

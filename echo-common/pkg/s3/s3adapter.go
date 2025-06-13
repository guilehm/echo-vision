package s3

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
	"github.com/rotisserie/eris"
)

type S3Adapter struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

func NewS3Adapter(bucketName, region string) (filestorage.FileStoragePort, error) {
	if bucketName == "" {
		return nil, eris.New("bucket name cannot be empty")
	}
	if region == "" {
		return nil, eris.New("region cannot be empty")
	}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
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

	s3Client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(s3Client)

	return &S3Adapter{
		client:    s3Client,
		presigner: presigner,
		bucket:    bucketName,
	}, nil
}

func (s *S3Adapter) GeneratePreSignedURL(fileKey filestorage.FileKey, contentType string) (string, error) {
	req, err := s.presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileKey.Filepath),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", eris.Wrap(err, "error generating pre-signed URL")
	}
	return req.URL, nil
}

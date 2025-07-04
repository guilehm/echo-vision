package filestoragemocks

import (
	"fmt"

	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
)

type FileStorageMock struct {
	bucket string
}

// GenerateFileURL implements filestorage.FileStoragePort.
func (m *FileStorageMock) GenerateFileURL(fileKey filestorage.FileKey) (string, error) {
	return fmt.Sprintf("https://s3-mock/%s/%s", m.bucket, fileKey.Filepath), nil
}

func NewFileStorageMock(bucketName string) *FileStorageMock {
	return &FileStorageMock{
		bucket: bucketName,
	}
}

func (m *FileStorageMock) GeneratePreSignedURL(fileKey filestorage.FileKey, contentType string) (string, error) {
	return fmt.Sprintf("https://s3-mock/%s/%s", m.bucket, fileKey.Filepath), nil
}

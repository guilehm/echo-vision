package filestorage

import "fmt"

type FileKey string

func (f FileKey) String() string {
	return string(f)
}

func NewFileKey(path, filename string) FileKey {
	return FileKey(fmt.Sprintf("%s/%s", path, filename))
}

type FileStoragePort interface {
	GeneratePreSignedURL(fileKey FileKey) (string, error)
}

type FileStorageMock struct {
	bucket string
}

func NewFileStorageMock(bucketName string) *FileStorageMock {
	return &FileStorageMock{
		bucket: bucketName,
	}
}

func (m *FileStorageMock) GeneratePreSignedURL(fileKey FileKey) (string, error) {
	return fmt.Sprintf("https://mock/%s/%s", m.bucket, fileKey.String()), nil
}

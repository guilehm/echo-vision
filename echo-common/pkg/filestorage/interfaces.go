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
	GeneratePreSignedURL(fileKey FileKey, contentType string) (string, error)
}

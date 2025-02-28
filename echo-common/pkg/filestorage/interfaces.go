package filestorage

import (
	"fmt"
	"strconv"
	"time"
)

type FileKey string

func (f FileKey) String() string {
	return string(f)
}

func NewFileKey(path, filename string) FileKey {
	return FileKey(fmt.Sprintf("%s/%s-%s", path, strconv.Itoa(int(time.Now().Unix())), filename))
}

type FileStoragePort interface {
	GeneratePreSignedURL(fileKey FileKey, contentType string) (string, error)
}

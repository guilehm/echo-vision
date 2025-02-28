package filestorage

import (
	"fmt"
	"strconv"
	"time"
)

type FileKey struct {
	// URL      string `json:"url"`
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

// func (f FileKey) String() string {
// 	return f.URL
// }

func NewFileKey(path, filename string) FileKey {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	return FileKey{
		// URL:      fmt.Sprintf("%s/%s-%s", path, timestamp, filename),
		Filename: filename,
		Filepath: fmt.Sprintf("%s/%s-%s", path, timestamp, filename),
	}
}

type FileStoragePort interface {
	GeneratePreSignedURL(fileKey FileKey, contentType string) (string, error)
}

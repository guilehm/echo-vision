package filestorage

import (
	"fmt"
	"strconv"
	"time"
)

type FileKey struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

func NewFileKey(path, filename string) FileKey {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	return FileKey{
		// URL:      fmt.Sprintf("%s/%s-%s", path, timestamp, filename),
		Filename: filename,
		Filepath: fmt.Sprintf("%s/%s-%s", path, timestamp, filename),
	}
}

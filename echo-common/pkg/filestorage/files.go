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

func NewUploadFileKey(path, filename string) FileKey {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	return FileKey{
		Filename: filename,
		Filepath: fmt.Sprintf("%s/%s-%s", path, timestamp, filename),
	}
}

func NewFileKey(filepath, filename string) FileKey {
	return FileKey{
		Filename: filename,
		Filepath: filepath,
	}
}

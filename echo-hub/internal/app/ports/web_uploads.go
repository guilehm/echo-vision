package ports

import (
	"net/http"
)

type UploadWebPort interface {
	PresignedURL(w http.ResponseWriter, r *http.Request)
}

type UploadPresignedURLInput struct {
	Filename    string `json:"filename"`
	Filepath    string `json:"filepath"`
	ContentType string `json:"contentType"`
}

type UploadPresignedURLResponse struct {
	URL string `json:"url"`
}

func (i UploadPresignedURLInput) IsValid() bool {
	return i.Filename != "" && i.Filepath != "" && i.ContentType != ""
}

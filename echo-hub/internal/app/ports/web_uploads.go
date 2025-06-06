package ports

import (
	"net/http"

	hubevents "github.com/guilehm/echo-vision/echo-hub/pkg/events"
)

type UploadWebPort interface {
	PresignedURL(w http.ResponseWriter, r *http.Request)
}

type UploadPresignedURLInput struct {
	Filename    string `json:"filename"`
	EventType   string `json:"eventType"`
	ContentType string `json:"contentType"`
}

type UploadPresignedURLResponse struct {
	URL      string `json:"url"`
	Filepath string `json:"filepath"`
	Filename string `json:"filename"`
}

func (i UploadPresignedURLInput) IsValid() bool {
	if i.Filename == "" || i.EventType == "" || i.ContentType == "" {
		return false
	}
	if !hubevents.EventType(i.EventType).IsValid() {
		return false
	}
	return true
}

package ports

import (
	"net/http"
)

type EventWebPort interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	ListEvents(w http.ResponseWriter, r *http.Request)
	ListOwnEvents(w http.ResponseWriter, r *http.Request)
	EventByID(w http.ResponseWriter, r *http.Request)
}

type EventCreateInput struct {
	EventType string `json:"eventType"`
	SubType   string `json:"subType"`

	Filename    string `json:"filename"`
	Filepath    string `json:"filepath"`
	ContentType string `json:"contentType"`
	Filesize    int64  `json:"filesize"`
}

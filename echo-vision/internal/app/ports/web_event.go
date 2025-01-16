package ports

import (
	"net/http"
)

type EventWebPort interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	ListEvents(w http.ResponseWriter, r *http.Request)
}

type EventCreateInput struct {
	EventType string `json:"eventType"`
	SubType   string `json:"subType"`
}

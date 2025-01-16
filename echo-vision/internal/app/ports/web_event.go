package ports

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type EventWebPort interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	ListEvents(w http.ResponseWriter, r *http.Request)
}

type EventCreateInput struct {
	EventType string `json:"eventType"`
	SubType   string `json:"subType"`
}

type EventResponse struct {
	UserID    uuid.UUID `json:"userID"`
	ID        uuid.UUID `json:"id"`
	EventType string    `json:"eventType"`
	SubType   string    `json:"subType"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updatedAt"`
}

type EventCreateResponse struct {
	ID uuid.UUID `json:"id"`
}

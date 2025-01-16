package web

import (
	"net/http"

	"github.com/guilehm/echo-vision/internal/app/ports"
)

type EventHandler struct {
	eventPort ports.EventPort
}

func NewEventHandler(up ports.EventPort) ports.EventWebPort {
	return &EventHandler{
		eventPort: up,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	panic("uninplemented")
}

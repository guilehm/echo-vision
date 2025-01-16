package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/infra/logging"
)

var logger = logging.NewLogger()

func NewRouter(up ports.UserPort, ep ports.EventPort) http.Handler {
	// handlers
	uh := NewUserHandler(up)
	eh := NewEventHandler(ep)

	// router
	r := chi.NewRouter()

	// middlewares
	r.Use(SetHeaders)
	r.Use(CorsMiddleware)
	r.Use(LogRequest)

	// routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.CreateUser)
	})
	r.Route("/events", func(r chi.Router) {
		r.Post("/", eh.CreateEvent)
	})
	return r
}

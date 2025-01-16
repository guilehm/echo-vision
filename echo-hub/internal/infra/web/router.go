package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/logging"
)

var logger = logging.NewLogger()

func NewRouter(up ports.UserPort, ep ports.EventPort) http.Handler {
	// handlers
	uh := NewUserHandler(up)
	eh := NewEventHandler(ep)

	// router
	r := chi.NewRouter()

	// middlewares
	authMiddleware := newAuthenticationMiddleware(up)

	r.Use(setHeaders)
	r.Use(corsMiddleware)
	r.Use(logRequest)

	// routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", uh.CreateUser)
		r.Post("/login", uh.Login)
		r.Post("/refresh-token", uh.RefreshToken)

		r.Route("/{userID}", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/events", eh.ListEvents)
		})
		r.Route("/me", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", uh.MeUser)
		})
	})

	r.Route("/events", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", eh.CreateEvent)
	})

	return r
}

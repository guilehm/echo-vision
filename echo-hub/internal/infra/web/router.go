package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/guilehm/echo-vision/echo-common/logging"
	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
)

var logger = logging.NewLogger()

func NewRouter(
	up ports.UserPort,
	ep ports.EventPort,
	upp filestorage.FileStoragePort,
	publisher messaging.Publisher,
) http.Handler {
	// handlers
	uh := NewUserHandler(up)
	eh := NewEventHandler(ep)
	uph := NewUploadHandler(upp)

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
			r.With(paginationMiddleware).Get("/events", eh.ListEvents)
		})
		r.Route("/me", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", uh.MeUser)
		})
	})

	r.Route("/events", func(r chi.Router) {
		r.Use(authMiddleware)
		r.With(paginationMiddleware).Get("/", eh.ListOwnEvents)
		r.Post("/", eh.CreateEvent)
	})

	r.Route("/uploads", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/presigned-url", uph.PresignedURL)
	})

	return r
}

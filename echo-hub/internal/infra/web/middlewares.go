package web

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/cors"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	"github.com/rotisserie/eris"
)

func corsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://echo-vision.local"}, // TODO: do not allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(next)
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(
			r.Context(),
			"request received",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("userAgent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}

func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func newAuthenticationMiddleware(userPort ports.UserPort) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				handleApiResponse(w, apiResponse[any](nil, newApiError(
					http.StatusForbidden,
					http.StatusText(http.StatusForbidden),
				)))
				return
			}

			ctx := r.Context()

			user, err := userPort.UserByAccessToken(ctx, token)
			if err != nil {
				handleApiResponse(w, apiResponse[any](nil, eris.Wrap(newApiError(
					http.StatusUnauthorized,
					shared.ErrInvalidToken.Error(),
				), shared.ErrInvalidID.Error())))
				return
			}

			slog.Debug("user authenticated", slog.String("email", user.Email()))
			ctx = context.WithValue(r.Context(), contextKeyMeUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func paginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := defaultPaginationParams()
		cursor := r.URL.Query().Get("cursor")
		limitQuery := r.URL.Query().Get("limit")

		limit := params.limit

		if limitQuery != "" {
			var err error
			limit, err = strconv.Atoi(limitQuery)
			if err != nil || limit <= 0 || limit > 100 {
				logger.Error("invalid limit parameter", slog.String("limit", limitQuery))
				handleApiResponse(w, apiResponse[any](nil, newApiError(
					http.StatusBadRequest,
					shared.ErrInvalidLimit.Error(),
				)))
				return
			}
		}

		params.limit = limit
		params.cursor = cursor

		ctx := context.WithValue(r.Context(), contextKeyPaginationParams, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

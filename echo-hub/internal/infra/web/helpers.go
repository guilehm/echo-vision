package web

import (
	"log/slog"
	"net/http"

	"github.com/lib/pq"
)

func logPGError(pgErr *pq.Error) (int, string) {
	logger.Error("database error",
		slog.String("error", pgErr.Error()),
		slog.String("code", string(pgErr.Code)),
		slog.String("detail", pgErr.Detail),
		slog.String("message", pgErr.Message),
	)
	switch pgErr.Code {
	case "23503":
		return http.StatusBadRequest, "foreign key violation constraint"
	case "23505":
		return http.StatusBadRequest, "unique constraint violation"
	case "22001":
		return http.StatusBadRequest, "string data right truncation"
	default:
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
}

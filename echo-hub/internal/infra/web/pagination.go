package web

import (
	"context"
	"log/slog"
)

type paginationParams struct {
	cursor string
	limit  int
}

var defaultPaginationParams = func() *paginationParams {
	return &paginationParams{
		cursor: "",
		limit:  10,
	}
}

func paginationParamsFromContext(ctx context.Context) *paginationParams {
	params, err := fromContext[paginationParams](ctx, contextKeyPaginationParams)
	if err != nil {
		logger.Error("error getting pagination params from context, returning default value", slog.String("error", err.Error()))
		return defaultPaginationParams()
	}
	return params
}

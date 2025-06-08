package web

import (
	"context"
	"log/slog"

	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
)

type contextKey string

const (
	contextKeyMeUser           contextKey = "me_user"
	contextKeyPaginationParams contextKey = "pagination_params"
)

func fromContext[T any](ctx context.Context, key contextKey) (*T, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, shared.ErrContextValueNotFound
	}

	t, ok := value.(*T)
	if !ok {
		return nil, shared.ErrContextValueNotFound
	}
	return t, nil
}

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

package web

import (
	"context"

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

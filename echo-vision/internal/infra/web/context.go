package web

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/shared"
)

type contextKey string

const (
	contextKeyMeUser contextKey = "me_user"
)

func fromContext[target any](ctx context.Context, key contextKey, t *target) (*target, error) {
	t, ok := ctx.Value(key).(*target)
	if !ok {
		return nil, shared.ErrContextValueNotFound
	}
	return t, nil
}

package repositories

import (
	"context"

	"github.com/guilehm/echo-vision/internal/app/domain"
)

type EventRepository interface {
	Save(ctx context.Context, tx Transaction, event *domain.Event) error

	BeginTx(ctx context.Context) (Transaction, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

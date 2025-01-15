package repositories

import "context"

type GenericRepository interface {
	BeginTx(ctx context.Context) (Transaction, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
}

type Transaction interface {
	Commit() error
	Rollback() error
}

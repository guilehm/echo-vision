package postgres

import (
	"context"
	"log/slog"
	"strings"

	"github.com/guilehm/echo-vision/internal/app/domain"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/logging"
	"github.com/rotisserie/eris"
)

var logger = logging.NewLogger()

func (r *EventRepository) BeginTx(ctx context.Context) (repositories.Transaction, error) {
	tx, err := r.entClient.BeginTx(ctx, nil)
	if err != nil {
		return nil, eris.Wrap(err, "error starting transaction")
	}
	return &entTransaction{tx: tx}, nil
}

// Save implements repositories.EventRepository.
func (e *EventRepository) Save(ctx context.Context, tx repositories.Transaction, event *domain.Event) error {
	panic("unimplemented")
}

// WithTransaction implements repositories.EventRepository.
func (r *EventRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx repositories.Transaction) error) error {
	logger.Debug("starting transaction")

	tx, err := r.BeginTx(ctx)
	if err != nil {
		return eris.Wrap(err, "error starting transaction")
	}
	defer func() {
		if v := recover(); v != nil {
			err := tx.Rollback()
			if err != nil {
				logger.Error("error rolling back transaction", slog.String("error", err.Error()))
			}
			panic(v)
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		logger.Debug("rolling back transaction", slog.String("error", err.Error()))
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Error("error rolling back transaction", slog.String("error", rollbackErr.Error()))
			return eris.Wrap(rollbackErr, "error rolling back transaction")
		}

		if strings.Contains(err.Error(), "could not obtain lock on row in relation") {
			logger.Warn("could not obtain lock on row", slog.String("error", err.Error()))
			return ErrCouldNotAcquireLock
		}
		return err
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		return eris.Wrap(commitErr, "error committing transaction")
	}

	logger.Debug("transaction successfully committed")
	return nil
}

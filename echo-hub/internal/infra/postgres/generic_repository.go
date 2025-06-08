package postgres

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/shared"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/rotisserie/eris"
)

// BeginTx implements repositories.Repository.
func (r *Repository) BeginTx(ctx context.Context) (repositories.Transaction, error) {
	tx, err := r.entClient.BeginTx(ctx, nil)
	if err != nil {
		return nil, eris.Wrap(err, "error starting transaction")
	}
	return &entTransaction{tx: tx}, nil
}

// WithTransaction implements repositories.Repository.
func (r *Repository) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx repositories.Transaction) error) error {
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

// resolveClient resolves the ent client
func (r *Repository) resolveClient(tx repositories.Transaction) *ent.Client {
	if tx == nil {
		return r.entClient
	}
	return tx.(*entTransaction).tx.Client()
}

// encodeCursor encodes createdAt|id into a base64 string
func encodeCursor(createdAt time.Time, id uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(
		fmt.Appendf(nil, "%s|%s", createdAt.UTC().Format(time.RFC3339Nano), id.String()),
	)
}

// decodeCursor decodes a base64 encoded cursor string into createdAt time and id
func decodeCursor(cursor string) (time.Time, uuid.UUID, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}
	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, uuid.UUID{}, shared.ErrInvalidCursor
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}
	id, err := uuid.Parse(parts[1])
	return t, id, err
}

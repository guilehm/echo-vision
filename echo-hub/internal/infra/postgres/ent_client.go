package postgres

import (
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	_ "github.com/lib/pq"
)

func NewEnt(dbURL string) *ent.Client {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	entClient := ent.NewClient(ent.Driver(drv))
	return entClient
}

type entTransaction struct {
	tx *ent.Tx
}

func (t *entTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *entTransaction) Rollback() error {
	return t.tx.Rollback()
}

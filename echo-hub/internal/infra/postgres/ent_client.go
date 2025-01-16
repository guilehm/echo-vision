package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	_ "github.com/lib/pq"
)

func NewEnt(dbURL string, schema string) *ent.Client {
	if dbURL == "" {
		log.Fatalf("missing DATABASE_URL")
	}
	if schema == "" {
		log.Fatalf("missing DATABASE_SCHEMA")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema))
	if err != nil {
		log.Fatalf("failed creating schema: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("SET search_path TO %s;", schema))
	if err != nil {
		log.Fatalf("failed setting search_path: %v", err)
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

//go:build exclude

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent/migrate"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		log.Fatal("migration name is required")
	}

	migrationName := os.Args[1]

	dbName := "postgres-migration"
	dbUser := "service-migrator"
	dbPassword := "password"
	var dbPort uint32 = 54322

	dbURL := fmt.Sprintf("postgresql://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, dbPassword, dbPort, dbName)
	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Password(dbPassword).
			Username(dbUser).
			Database(dbName).
			Version(embeddedpostgres.V12).
			Port(dbPort),
	)
	err := postgres.Start()
	if err != nil {
		log.Fatalf("failed starting embedded postgres: %v", err)
	}
	defer postgres.Stop()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occurred:", r)
			postgres.Stop()
		}
	}()

	dir, err := sqltool.NewGolangMigrateDir("internal/infra/postgres/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	opts := []schema.MigrateOption{
		schema.WithFormatter(sqltool.GolangMigrateFormatter),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.Postgres),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
		schema.WithDir(dir),
	}

	err = migrate.NamedDiff(ctx, dbURL, migrationName, opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}

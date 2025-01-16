package tests

import (
	"fmt"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/infra/postgres"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	pg        *embeddedpostgres.EmbeddedPostgres
	entClient *ent.Client
	repo      repositories.Repository
	m         *migrate.Migrate
)

var _ = BeforeSuite(func() {
	dbName := "echo-vision"
	dbUser := "echo-vision"
	dbPassword := "echo-vision"
	var dbPort uint32 = 15432

	dbURL := fmt.Sprintf("postgresql://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, dbPassword, dbPort, dbName)
	pg = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Password(dbPassword).
		Username(dbUser).
		Database(dbName).
		Version(embeddedpostgres.V15).
		Port(dbPort),
	)
	err := pg.Start()
	Expect(err).ToNot(HaveOccurred())
	entClient = postgres.NewEnt(dbURL)
	repo = postgres.NewRepository(entClient)

	migrationsDir := "file://../internal/infra/postgres/migrations"
	m, err = migrate.New(migrationsDir, dbURL)
	Expect(err).ToNot(HaveOccurred())
})

func TestEchoVision(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EchoVision Suite")
}

var _ = BeforeEach(func() {
	Expect(m.Up()).To(Succeed())
})

var _ = AfterEach(func() {
	Expect(m.Down()).To(Succeed())
})

var _ = AfterSuite(func() {
	_ = pg.Stop()
})

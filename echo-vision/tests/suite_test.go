package tests

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilehm/echo-vision/internal/app/ports"
	"github.com/guilehm/echo-vision/internal/app/repositories"
	"github.com/guilehm/echo-vision/internal/app/usecases"
	bcrypthasher "github.com/guilehm/echo-vision/internal/infra/bcrypt_hasher"
	jwtadapter "github.com/guilehm/echo-vision/internal/infra/jwt"
	"github.com/guilehm/echo-vision/internal/infra/postgres"
	"github.com/guilehm/echo-vision/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/internal/infra/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	pg              *embeddedpostgres.EmbeddedPostgres
	entClient       *ent.Client
	repo            repositories.Repository
	m               *migrate.Migrate
	server          *httptest.Server
	passwordAdapter ports.PasswordManager
)

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

var _ = BeforeSuite(func() {
	dbName := "echo-vision"
	dbUser := "echo-vision"
	dbPassword := "echo-vision"
	var dbPort uint32 = 15432

	// setup embedded postgres
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

	// setup ent postgres
	entClient = postgres.NewEnt(dbURL)
	repo = postgres.NewRepository(entClient)

	// run migrations
	migrationsDir := "file://../internal/infra/postgres/migrations"
	m, err = migrate.New(migrationsDir, dbURL)
	Expect(err).ToNot(HaveOccurred())

	// setup token manager
	jwtAdapter := jwtadapter.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		1*time.Hour,
		24*time.Hour,
	)

	// setup bcrypt password manager
	passwordAdapter = bcrypthasher.NewBcryptAdapter()

	// setup usecases
	userUseCase := usecases.NewManageUsersUseCase(repo, jwtAdapter, passwordAdapter)
	eventUseCase := usecases.NewManageEventsUseCase(repo)

	// setup http server
	router := web.NewRouter(userUseCase, eventUseCase)
	server = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	_ = pg.Stop()
	server.Close()
})

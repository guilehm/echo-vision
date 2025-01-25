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
	rabbitmqmocks "github.com/guilehm/echo-vision/echo-common/rabbitmq/mocks"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/usecases"
	bcrypthasher "github.com/guilehm/echo-vision/echo-hub/internal/infra/bcrypt_hasher"
	jwtadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/jwt"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/web"
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
	userUseCase     ports.UserPort
	eventUseCase    ports.EventPort
	jwtSecretKey    = os.Getenv("JWT_SECRET")
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
	dbSchema := "echo_hub"
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%d/%s?sslmode=disable&search_path=%s",
		dbUser,
		dbPassword,
		dbPort,
		dbName,
		dbSchema,
	)
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
	entClient = postgres.NewEnt(dbURL, dbSchema)
	repo = postgres.NewRepository(entClient)

	// run migrations
	migrationsDir := "file://../internal/infra/postgres/migrations"
	m, err = migrate.New(migrationsDir, dbURL)
	Expect(err).ToNot(HaveOccurred())

	// setup token manager
	jwtAdapter := jwtadapter.NewJWTManager(
		jwtSecretKey,
		1*time.Hour,
		24*time.Hour,
	)

	// setup bcrypt password manager
	passwordAdapter = bcrypthasher.NewBcryptAdapter()

	// setup rabbitmq mocks
	publisher := rabbitmqmocks.NewPublisher()

	// setup usecases
	userUseCase = usecases.NewManageUsersUseCase(repo, jwtAdapter, passwordAdapter)
	eventUseCase = usecases.NewManageEventsUseCase(repo, publisher)

	// setup http server
	router := web.NewRouter(userUseCase, eventUseCase, publisher)
	server = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	_ = pg.Stop()
	server.Close()
})

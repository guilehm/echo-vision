package tests

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilehm/echo-vision/echo-common/pkg/filestorage"

	filestoragemocks "github.com/guilehm/echo-vision/echo-common/pkg/filestorage/mocks"
	"github.com/guilehm/echo-vision/echo-common/pkg/messaging"
	rabbitmqmocks "github.com/guilehm/echo-vision/echo-common/rabbitmq/mocks"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/ports"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/repositories"
	"github.com/guilehm/echo-vision/echo-hub/internal/app/usecases"
	bcrypthasher "github.com/guilehm/echo-vision/echo-hub/internal/infra/bcrypt_hasher"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/consumers"
	jwtadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/jwt"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres"
	"github.com/guilehm/echo-vision/echo-hub/internal/infra/postgres/generated/ent"

	rabbitmqadapter "github.com/guilehm/echo-vision/echo-hub/internal/infra/rabbitmq"
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
	s3Mock          filestorage.FileStoragePort
	jwtSecretKey    = os.Getenv("JWT_SECRET")
	publisher       messaging.Publisher
	consumer        messaging.Consumer
	adapter         *rabbitmqadapter.RabbitMQAdapter
	handler         *rabbitmqmocks.Handler
)

func TestEchoHub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EchoHub Suite")
}

var _ = BeforeEach(func() {
	Expect(m.Up()).To(Succeed())
	t := GinkgoT()
	handler.Mock.Test(t)
	DeferCleanup(func() {
		Expect(handler.Mock.AssertExpectations(t)).To(BeTrue())
		handler.ResetMock()
	})
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

	// setup s3 mock
	s3Mock = filestoragemocks.NewFileStorageMock("echo-hub")

	// setup rabbitmq mocks
	mockChan := make(chan messaging.Message)
	publisher = rabbitmqmocks.NewPublisher(mockChan)

	// setup usecases
	userUseCase = usecases.NewManageUsersUseCase(repo, jwtAdapter, passwordAdapter)
	eventUseCase = usecases.NewManageEventsUseCase(repo, publisher)

	consumerGroup := consumers.NewConsumerGroup(eventUseCase)
	consumer = rabbitmqmocks.NewConsumer(mockChan)
	adapter = rabbitmqadapter.NewRabbitMQAdapter(consumerGroup)
	handler = rabbitmqmocks.NewHandler()

	go consumer.Subscribe(context.Background(), handler)
	go publisher.StartPublisher(context.Background())

	// setup http server
	router := web.NewRouter(userUseCase, eventUseCase, s3Mock, publisher)
	server = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	err := pg.Stop()
	Expect(err).ToNot(HaveOccurred())
	server.Close()
})

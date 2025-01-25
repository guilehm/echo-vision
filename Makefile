DOCKER_COMPOSE = docker-compose
MIGRATIONS_PATH = ./$(SERVICE_NAME)/internal/infra/postgres/migrations

generate_template:
	@echo "running generate for $(SERVICE_NAME)..."
	@cd $(SERVICE_NAME) && go run ./internal/infra/postgres/generated/entc.go

migrate_template:
	@echo "applying migrations for: $(SERVICE_NAME)"
	migrate -path $(MIGRATIONS_PATH) \
		-database "postgres://postgres:postgres@localhost:5432/echo-vision?sslmode=disable&search_path=$(SCHEMA_NAME)" up

run:
	@echo "starting all containers"
	$(DOCKER_COMPOSE) up -d

run_hub:
	@echo "starting echo-hub"
	$(DOCKER_COMPOSE) up echo-hub

setup:
	@echo "starting rabbitmq and postgres"
	$(DOCKER_COMPOSE) up -d rabbitmq postgres

stop:
	@echo "stopping all containers"
	$(DOCKER_COMPOSE) stop

down:
	@echo "stopping all containers"
	$(DOCKER_COMPOSE) down

remove:
	@echo "stopping all containers and removing volumes"
	$(DOCKER_COMPOSE) down -v

test_hub:
	@echo "running tests for echo-hub"
	TZ=UTC ginkgo -v echo-hub/tests

test: test_hub
	@echo "all tests completed"

migrate:
	@$(MAKE) migrate_template SERVICE_NAME=echo-hub SCHEMA_NAME=echo_hub

generate:
	@$(MAKE) generate_template SERVICE_NAME=echo-hub

# VARIABLES
DOCKER_COMPOSE = docker-compose
MIGRATIONS_PATH = ./$(SERVICE_NAME)/internal/infra/postgres/migrations

# TEMPLATES
generate_template:
	@echo "running generate for $(SERVICE_NAME)..."
	@cd $(SERVICE_NAME) && go run ./internal/infra/postgres/generated/entc.go

migrate_template:
	@echo "applying migrations for: $(SERVICE_NAME)"
	migrate -path $(MIGRATIONS_PATH) \
		-database "postgres://postgres:postgres@localhost:5432/echo-vision?sslmode=disable&search_path=$(SCHEMA_NAME)" up

# GLOBAL
setup:
	@echo "starting rabbitmq and postgres"
	$(DOCKER_COMPOSE) up -d rabbitmq postgres

run:
	@echo "starting all containers"
	$(DOCKER_COMPOSE) up -d

stop:
	@echo "stopping all containers"
	$(DOCKER_COMPOSE) stop

down:
	@echo "stopping all containers"
	$(DOCKER_COMPOSE) down

remove:
	@echo "stopping all containers and removing volumes"
	$(DOCKER_COMPOSE) down -v

test: clear_database test_hub test_analyzer
	@echo "all tests completed"

migrate:
	@$(MAKE) migrate_template SERVICE_NAME=echo-hub SCHEMA_NAME=echo_hub
	# @$(MAKE) migrate_template SERVICE_NAME=echo-analyzer SCHEMA_NAME=echo_analyzer

generate:
	@$(MAKE) generate_template SERVICE_NAME=echo-hub

# SERVICES
run_analyzer:
	@echo "starting echo-analyzer"
	$(DOCKER_COMPOSE) up echo-analyzer

run_hub:
	@echo "starting echo-hub"
	$(DOCKER_COMPOSE) up echo-hub

test_analyzer: clear_database
	@echo "running tests for echo-analyzer"
	TZ=UTC ginkgo -v --race echo-analyzer/tests

test_hub: clear_database
	@echo "running tests for echo-hub"
	TZ=UTC ginkgo -v --race echo-hub/tests

clear_database:
	lsof -ti tcp:15432 | xargs kill -9 || true

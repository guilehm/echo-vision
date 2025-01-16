# Echo-Vision Project

## Overview

The `echo-vision` project provides services for user authentication and interaction endpoints, managed via microservices. This project uses Docker for containerization, RabbitMQ for messaging, and PostgreSQL for database operations. Below are the available commands in the `Makefile` to streamline development tasks.

---

## Prerequisites

- **Docker**: Ensure Docker is installed and running.
- **Go**: Used for code generation and application development.
- **Migrate**: For database migration operations.
- **Ginkgo**: Test framework for running tests.

---

## Commands

### General Commands

- **Setup essential services**:
  ```bash
  make setup
  ```
  Starts `RabbitMQ` and `PostgreSQL` containers.

- **Run all containers**:
  ```bash
  make run
  ```
  Starts all Docker containers for the project.

- **Run a specific service (echo-hub)**:
  ```bash
  make run_hub
  ```
  Starts only the `echo-hub` container.

- **Stop all containers**:
  ```bash
  make stop
  ```
  Stops all running containers.

- **Remove containers and volumes**:
  ```bash
  make remove
  ```
  Stops all containers and removes associated volumes.

---

### Code Generation and Migrations

- **Generate database templates**:
  ```bash
  make generate
  ```
  Runs code generation for `echo-hub` service's database templates.

- **Apply database migrations**:
  ```bash
  make migrate
  ```
  Applies migrations for the `echo-hub` service using the `migrate` tool.

---

### Testing

- **Run tests for `echo-hub`**:
  ```bash
  make test_hub
  ```
  Executes tests for the `echo-hub` service using Ginkgo.

- **Run all tests**:
  ```bash
  make test_all
  ```
  Runs all test suites across the project.

---

## Directory Structure

- `docker-compose.yml`: Contains the service definitions for Docker containers.
- `./echo-hub/internal/infra/postgres/migrations`: Path for database migrations for `echo-hub`.
- `./echo-hub/internal/infra/postgres/generated/entc.go`: Script for database template generation.

---

## Notes

1. Replace placeholders such as `SERVICE_NAME` and `SCHEMA_NAME` when needed. Default service is `echo-hub` and schema is `echo_hub`.
2. Ensure all required tools (e.g., `migrate`, `ginkgo`) are installed and available in your environment.

---

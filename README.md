# Echo-Vision Project

## Overview

Echo-Vision is a distributed system composed of microservices designed to analyze images and videos using AI, manage user authentication, and expose interaction endpoints through a modern dashboard. The platform uses:

- **Docker** for containerized environments.
- **RabbitMQ** for asynchronous communication between services.
- **PostgreSQL** for relational data persistence.
- **Go** with Clean Architecture + DDD for service implementation.
- **Hosts Configuration**: Add `echo-vision.local` to `/etc/hosts` for proper service resolution.

Each component plays a specific role within the architecture:

- **echo-hub**: Manages user authentication, session handling, and provides endpoints for the frontend to create and manage events.
- **echo-analyzer**: Processes image and video files using AWS Rekognition, extracting metadata such as labels, faces, and objects, and returns structured analysis results through RabbitMQ.
- **echo-common**: A shared Go module (not a standalone service) containing reusable packages such as DTOs, message definitions, constants, and utility functions. It ensures consistency and reduces duplication across all microservices.
- **echo-ui**: The frontend application built with Next.js and React, providing an intuitive dashboard for users to upload media, manage their sessions, view analysis results, and interact with the system in real time.

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
  make test
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

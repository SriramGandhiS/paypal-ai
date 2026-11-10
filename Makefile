# ==============================================================================
# PayBridge AI — Root Developer Makefile
# ==============================================================================
# Usage:
#   make setup          - Verify environment and tools
#   make up             - Start local Docker infrastructure (Postgres, Redis, Mailpit)
#   make down           - Stop local Docker infrastructure
#   make status         - Inspect running container health
#   make test           - Run backend and contract tests
#   make lint           - Run Go vet and static checks
#   make codegen        - Generate Go and TypeScript types from OpenAPI 3.1
#   make migrate-up     - Run pending database migrations
#   make migrate-status - Inspect current database migration version
# ==============================================================================

.PHONY: help setup up down status test lint codegen migrate-up migrate-status clean

DOCKER_COMPOSE_FILE = infra/docker/docker-compose.yml
DB_URL ?= postgres://paybridge_admin:paybridge_dev_secret@localhost:5432/paybridge_dev?sslmode=disable

help:
	@echo "PayBridge AI Development Commands:"
	@echo "  make setup          - Verify tools and environment prerequisites"
	@echo "  make up             - Start local Docker infrastructure"
	@echo "  make down           - Stop local Docker infrastructure"
	@echo "  make status         - Check status of local containers"
	@echo "  make test           - Run backend unit and integration tests"
	@echo "  make lint           - Run Go code vet and linting"
	@echo "  make codegen        - Regenerate Go and TS types from openapi.yaml"
	@echo "  make migrate-up     - Run goose database migrations"
	@echo "  make migrate-status - Check goose database migration status"

setup:
	@echo "==> Verifying environment prerequisites..."
	@which go > /dev/null 2>&1 && go version || echo "WARNING: Go not found in PATH"
	@which node > /dev/null 2>&1 && node -v || echo "WARNING: Node.js not found in PATH"
	@which docker > /dev/null 2>&1 && docker version || echo "WARNING: Docker not found in PATH"
	@test -f .env || (echo "==> Creating .env from .env.example..." && cp .env.example .env)
	@echo "==> Setup verification complete."

up:
	@echo "==> Starting PayBridge Docker infrastructure..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "==> Waiting for services to become healthy..."
	docker compose -f $(DOCKER_COMPOSE_FILE) ps

down:
	@echo "==> Stopping PayBridge Docker infrastructure..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

status:
	@echo "==> PayBridge Container Status:"
	docker compose -f $(DOCKER_COMPOSE_FILE) ps

test:
	@echo "==> Running Go backend tests with race detector..."
	@if [ -f services/api/go.mod ]; then \
		cd services/api && go test -v -race ./...; \
	else \
		echo "services/api/go.mod not yet initialized."; \
	fi

lint:
	@echo "==> Running Go vet..."
	@if [ -f services/api/go.mod ]; then \
		cd services/api && go vet ./...; \
	fi

codegen:
	@echo "==> Running OpenAPI code generation..."
	@if [ -f scripts/codegen.sh ]; then \
		bash scripts/codegen.sh; \
	else \
		echo "scripts/codegen.sh not found."; \
	fi

migrate-up:
	@echo "==> Running database migrations..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir services/api/db/migrations postgres "$(DB_URL)" up; \
	else \
		go run ./services/api/cmd/migrate -action=up -db-url="$(DB_URL)"; \
	fi

migrate-status:
	@echo "==> Checking database migration status..."
	@if command -v goose > /dev/null 2>&1; then \
		goose -dir services/api/db/migrations postgres "$(DB_URL)" status; \
	else \
		go run ./services/api/cmd/migrate -action=status -db-url="$(DB_URL)"; \
	fi

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/ dist/ apps/web/dist/

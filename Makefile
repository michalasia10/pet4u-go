# Simple Makefile for a Go project

# Tools
GINKGO := $(shell go env GOPATH)/bin/ginkgo
MIGRATIONS_DIR := ./migrations

# --- Env loading ---
# Load environment variables from .env automatically for CLI targets
ENV_FILE := .env
LOAD_ENV := set -a; if [ -f $(ENV_FILE) ]; then . $(ENV_FILE); fi; set +a

# Build the application
all: build-binary test

build-binary:
	@echo "Building binary..."
	@mkdir -p bin
	@go build -o bin/api ./cmd/api

# Run the application
run:
	@go run ./cmd/api

# --- Migrations (Goose) ---
# Generate new Go migration skeleton (requires name=...)
makemigrations:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) -s create $(name)

# Run migrations using our custom binary (imports Go migrations)
migrate:
	@$(LOAD_ENV); \
	if [ -z "$$BLUEPRINT_DB_USERNAME" ] || [ -z "$$BLUEPRINT_DB_PASSWORD" ] || [ -z "$$BLUEPRINT_DB_HOST" ] || [ -z "$$BLUEPRINT_DB_PORT" ] || [ -z "$$BLUEPRINT_DB_DATABASE" ] || [ -z "$$BLUEPRINT_DB_SCHEMA" ]; then \
		echo "Set BLUEPRINT_DB_* env vars before running migrate"; \
		exit 1; \
	fi
	@go run ./cmd/migrate up

migrate-down:
	@$(LOAD_ENV); \
	if [ -z "$$BLUEPRINT_DB_USERNAME" ] || [ -z "$$BLUEPRINT_DB_PASSWORD" ] || [ -z "$$BLUEPRINT_DB_HOST" ] || [ -z "$$BLUEPRINT_DB_PORT" ] || [ -z "$$BLUEPRINT_DB_DATABASE" ] || [ -z "$$BLUEPRINT_DB_SCHEMA" ]; then \
		echo "Set BLUEPRINT_DB_* env vars before running migrate-down"; \
		exit 1; \
	fi
	@go run ./cmd/migrate down

migrate-status:
	@$(LOAD_ENV); \
	if [ -z "$$BLUEPRINT_DB_USERNAME" ] || [ -z "$$BLUEPRINT_DB_PASSWORD" ] || [ -z "$$BLUEPRINT_DB_HOST" ] || [ -z "$$BLUEPRINT_DB_PORT" ] || [ -z "$$BLUEPRINT_DB_DATABASE" ] || [ -z "$$BLUEPRINT_DB_SCHEMA" ]; then \
		echo "Set BLUEPRINT_DB_* env vars before running migrate-status"; \
		exit 1; \
	fi
	@go run ./cmd/migrate status

dev-up:
	@docker compose -f docker-compose.dev.yml up -d

dev-rebuild:
	@docker compose -f docker-compose.dev.yml up -d --build


dev-down:
	@docker compose -f docker-compose.dev.yml down -v

dev-logs:
	@docker compose -f docker-compose.dev.yml logs -f api

# Test the application
test:
	@echo "Testing (fast specs)..."
	@$(GINKGO) -r -p -v --skip-package=internal/database

test-db:
	@echo "Testing (db)..."
	@$(GINKGO) -v ./internal/database

test-all:
	@echo "Testing (all specs)..."
	@make test
	@make test-db



# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/api

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build-binary run test clean watch dev-up dev-down dev-rebuild dev-logs itest migrate-create migrate-up migrate-down migrate-status

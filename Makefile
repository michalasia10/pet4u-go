# Simple Makefile for a Go project

# Tools
GINKGO := $(shell go env GOPATH)/bin/ginkgo

# Build the application
all: build-binary test

build-binary:
	@echo "Building binary..."
	@mkdir -p bin
	@go build -o bin/api ./cmd/api

# Run the application
run:
	@go run ./cmd/api

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

.PHONY: all build-binary run test clean watch dev-up dev-down dev-rebuild dev-logs itest

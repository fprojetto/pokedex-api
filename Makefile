APP_NAME := pokedex-api
APP_TAG := latest
BUILD_DIR := ./bin
GO := go
DOCKER := docker

.PHONY: all build docker-build docker-run run clean help

all: build

build:
	@echo ">> Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

test:
	@echo ">> Running tests for $(APP_NAME)..."
	@$(GO) test ./... -coverprofile=coverage.out
	@$(GO) tool cover -func=coverage.out
	@echo "✅ Running tests complete"

docker-build:
	@echo ">> Building docker image $(APP_NAME):$(APP_TAG)..."
	@DOCKER_BUILDKIT=1 $(DOCKER) build --no-cache -f Dockerfile -t $(APP_NAME):$(APP_TAG) .
	@echo "🐳 Docker image built"

docker-run:
	@echo ">> Running docker container $(APP_NAME):$(APP_TAG)..."
	@$(DOCKER) run -p 8080:8080 \
		$(APP_NAME):$(APP_TAG)

run:
	@echo ">> Running $(APP_NAME)..."
	@$(GO) run ./$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo ">> Cleaning $(APP_NAME) build artifacts..."
	@rm -f $(BUILD_DIR)/$(APP_NAME)
	@echo "🧹 Clean complete."

help:
	@echo ""
	@echo "Available commands:"
	@echo "  make build                 Build the app binary"
	@echo "  make docker-build          Build the docker container for the app"
	@echo "  make docker-run            Run the app using docker"
	@echo "  make test                  Run tests"
	@echo "  make run                   Run the app locally"
	@echo "  make clean                 Remove built app binary"
	@echo ""

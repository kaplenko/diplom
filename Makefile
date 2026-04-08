# ──────────────────────────────────────────────
#  Diplom — development Makefile
# ──────────────────────────────────────────────

APP_NAME   := diplom
CMD_DIR    := ./cmd/$(APP_NAME)
BUILD_DIR  := ./bin
BINARY     := $(BUILD_DIR)/$(APP_NAME)
SEED_CMD   := ./cmd/seed

.PHONY: run seed build test smoke lint tidy docker-up docker-down help

## run: start the main application
run:
	go run $(CMD_DIR)/main.go

## seed: populate the database with test data
seed:
	go run $(SEED_CMD)/main.go

## build: compile the Go binary into ./bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(CMD_DIR)
	@echo "Binary built: $(BINARY)"

## test: run all unit tests
test:
	go test ./...

## smoke: run smoke tests against a running server
smoke:
	@echo "Running smoke tests..."
	bash smoke-test.sh

## lint: run static analysis (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null 2>&1 \
		&& golangci-lint run ./... \
		|| echo "golangci-lint not installed — skipping (https://golangci-lint.run/welcome/install/)"

## tidy: clean up go.mod / go.sum
tidy:
	go mod tidy

## docker-up: start all services via docker-compose
docker-up:
	docker-compose up -d --build

## docker-down: stop all docker-compose services
docker-down:
	docker-compose down

## help: show available targets
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  make /' | sed 's/:/\t—/'
	@echo ""

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BINARY_NAME=learning-api
MIGRATIONS_DIR=internal/persistance/migrations

.PHONY: build run clean sqlc test swagger tidy migrate-up migrate-down migrate-status db-dump test-gcp-integration

tidy:
	@echo "Tidying go modules..."
	@go mod tidy

build:
	@echo "Building..."
	@go build -o tmp/$(BINARY_NAME) cmd/api/main.go

run: build
	@echo "Running..."
	@./tmp/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -rf tmp/

sqlc:
	@echo "Generating code with sqlc..."
	@go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

swagger:
	@echo "Generating swagger docs..."
	@go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go -o ./docs
	@echo "Converting to OpenAPI 3.0..."
	@go run internal/tools/convert_swagger.go

test:
	@echo "Running tests..."
	@go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-gcp-integration:
	@echo "Running GCP file service integration test..."
	@go test ./internal/gcp -run TestFileServiceUploadIntegration -v

# Postgres targets (Real DB)
migrate-up:
	@echo "Running migrations up (Postgres)..."
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	@echo "Running migrations down (Postgres)..."
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	@echo "Migration status (Postgres)..."
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

db-dump:
	@echo "Dumping schema to schema.sql..."
	@pg_dump --schema-only --no-owner "$(DB_URL)" > schema.sql

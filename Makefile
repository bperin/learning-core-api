ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BINARY_NAME=learning-api
MIGRATIONS_DIR=internal/persistance/migrations

.PHONY: build run clean sqlc test swagger tidy migrate-up migrate-down migrate-status db-dump

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

test:
	@echo "Running tests..."
	@go test ./...

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

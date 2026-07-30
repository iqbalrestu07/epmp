.PHONY: run build test lint tidy migrate-up migrate-down migrate-create

# ─── Development ──────────────────────────────────────────────────────────────

run:
	go run ./cmd/server/...

build:
	go build -o bin/server ./cmd/server/...

tidy:
	go mod tidy

# ─── Testing ──────────────────────────────────────────────────────────────────

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ─── Linting ──────────────────────────────────────────────────────────────────

lint:
	golangci-lint run ./...

# ─── Database Migrations ──────────────────────────────────────────────────────

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-version:
	go run ./cmd/migrate version

# Usage: make migrate-create name=add_column_to_tenants
migrate-create:
	go run ./cmd/migrate create $(name)

# Watch migrations/ and auto-apply on new files
migrate-watch:
	go run ./cmd/migrate watch

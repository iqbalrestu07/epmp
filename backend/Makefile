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
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

# Usage: make migrate-create name=create_properties
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

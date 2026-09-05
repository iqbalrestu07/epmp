.PHONY: all test test-e2e e2e test-backend test-frontend dev-backend dev-frontend build clean help

# Default target
all: help

help:
	@echo "=================================================================="
	@echo "               EPMP — Fullstack Management Commands               "
	@echo "=================================================================="
	@echo "  make test-e2e       : Run Playwright E2E tests (Google Chrome Headless)"
	@echo "  make test-e2e-gui   : Run Playwright E2E with visible Chrome GUI (Headed)"
	@echo "  make e2e            : Alias for test-e2e"
	@echo "  make e2e-gui        : Alias for test-e2e-gui"
	@echo "  make test-backend   : Run Go unit and integration tests"
	@echo "  make test-frontend  : Run TypeScript typecheck & Vite build"
	@echo "  make test           : Run full test suite (backend + E2E)"
	@echo "  make dev-backend    : Run Go Echo server (port 8080)"
	@echo "  make dev-frontend   : Run Vite frontend server (port 3000)"
	@echo "  make build          : Build both backend binary and frontend assets"
	@echo "=================================================================="

# ─── End-to-End Testing ───────────────────────────────────────────────────────

test-e2e:
	@echo "🚀 Running Playwright End-to-End Test Suite with System Chrome (Headless)..."
	cd frontend && npm run test:e2e

e2e: test-e2e

test-e2e-gui:
	@echo "🖥️ Launching Playwright E2E with visible Google Chrome GUI window..."
	cd frontend && npm run test:e2e:gui

e2e-gui: test-e2e-gui

# ─── Backend & Frontend Unit/Integration Tests ────────────────────────────────

test-backend:
	@echo "🧪 Running Backend Tests (Go)..."
	cd backend && go test ./... -v -count=1

test-frontend:
	@echo "🔬 Running Frontend TypeScript & Build Verification..."
	cd frontend && npm run build

test: test-backend test-frontend test-e2e
	@echo "🎉 All tests passed successfully!"

# ─── Development Servers ──────────────────────────────────────────────────────

dev-backend:
	@echo "⚡ Starting Backend Go Server (http://localhost:8080)..."
	cd backend && go run ./cmd/server/...

dev-frontend:
	@echo "⚡ Starting Frontend Vite Server (http://localhost:3000)..."
	cd frontend && npm run dev

# ─── Build ────────────────────────────────────────────────────────────────────

build:
	@echo "📦 Building Backend..."
	cd backend && go build -o bin/server ./cmd/server/...
	@echo "📦 Building Frontend..."
	cd frontend && npm run build
	@echo "✅ Build completed!"

clean:
	rm -rf backend/bin frontend/dist
	@echo "🧹 Clean completed!"

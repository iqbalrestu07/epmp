# ─────────────────────────────────────────────────────────────────────────────
# EPMP Backend — Dockerfile
# Multi-stage build:
#   Stage 1 (builder) — compile the Go binary inside the official Go image
#   Stage 2 (runner)  — copy only the binary into a minimal distroless image
#
# Build:
#   docker build -t epmp-backend:latest .
#
# Run:
#   docker run --env-file .env -p 8080:8080 epmp-backend:latest
# ─────────────────────────────────────────────────────────────────────────────

# ── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

# Install git (required by some go get / module resolution steps)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache module downloads separately from source so rebuilds are fast
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build a statically linked binary
#   CGO_ENABLED=0  — pure Go, no C deps
#   -trimpath      — remove local filesystem paths from binary
#   -ldflags       — strip debug info + DWARF to reduce image size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /build/bin/server \
      ./cmd/server/...

# Build the migration runner as a separate binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /build/bin/migrate \
      ./cmd/migrate/...

# ── Stage 2: Runner ───────────────────────────────────────────────────────────
# gcr.io/distroless/static-debian12 is a minimal image with:
#   - CA certificates
#   - timezone data
#   - No shell, no package manager (attack surface minimised)
FROM gcr.io/distroless/static-debian12 AS runner

# Run as non-root (distroless provides user 65532 "nonroot")
USER nonroot:nonroot

WORKDIR /app

# Copy binaries
COPY --from=builder --chown=nonroot:nonroot /build/bin/server  /app/server
COPY --from=builder --chown=nonroot:nonroot /build/bin/migrate /app/migrate

# Copy migration SQL files (runner reads them at startup)
COPY --from=builder --chown=nonroot:nonroot /build/migrations /app/migrations

# Copy timezone data + CA certs from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# Default: run the API server
# Override CMD to run migrations: docker run ... /app/migrate up
ENTRYPOINT ["/app/server"]

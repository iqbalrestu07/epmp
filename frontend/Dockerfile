# ─────────────────────────────────────────────────────────────────────────────
# EPMP Frontend — Dockerfile
# Multi-stage build:
#   Stage 1 (builder) — install deps and build Vite production bundle
#   Stage 2 (runner)  — serve /dist with a minimal nginx image
#
# Build:
#   docker build \
#     --build-arg VITE_API_BASE_URL=https://api.epmp.example.com \
#     -t epmp-frontend:latest .
#
# Run:
#   docker run -p 3000:80 epmp-frontend:latest
# ─────────────────────────────────────────────────────────────────────────────

# ── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM node:20-alpine AS builder

# Build-time API URL (baked into the bundle by Vite)
ARG VITE_API_BASE_URL=http://localhost:8080
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

WORKDIR /build

# Cache node_modules layer — only reinstall when package.json changes
COPY package.json package-lock.json* ./
RUN npm ci --prefer-offline

# Copy source and build
COPY . .
RUN npm run build

# ── Stage 2: Runner ───────────────────────────────────────────────────────────
FROM nginx:1.27-alpine AS runner

# Remove default nginx page
RUN rm -rf /usr/share/nginx/html/*

# Copy built assets
COPY --from=builder /build/dist /usr/share/nginx/html

# Custom nginx config for SPA (handles client-side routing)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

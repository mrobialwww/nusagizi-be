# =============================================================================
# Stage 1: Base — shared dependencies
# =============================================================================
FROM golang:1.25-alpine AS base

RUN apk add --no-cache git ca-certificates tzdata curl netcat-openbsd

WORKDIR /app

# Copy module files first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Install migrate CLI — pinned version for reproducible builds
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# =============================================================================
# Stage 2: Development — pre-built binary with debug info
# =============================================================================
FROM base AS dev

WORKDIR /app

COPY . .

# Build binary with debug info (no -w -s flags unlike prod)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/server ./cmd/api/main.go

RUN chmod +x ./scripts/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./scripts/entrypoint.sh"]

# =============================================================================
# Stage 3: Builder — compile static binary
# =============================================================================
FROM base AS builder

WORKDIR /app

COPY . .

# CGO disabled for a fully static binary; strip debug info to reduce size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

# =============================================================================
# Stage 4: Production — minimal runtime image
# =============================================================================
FROM alpine:3.20 AS prod

RUN apk add --no-cache ca-certificates tzdata curl netcat-openbsd

ENV TZ=Asia/Jakarta

# Run as non-root for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy only required artifacts from builder
COPY --from=builder /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

USER appuser

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]

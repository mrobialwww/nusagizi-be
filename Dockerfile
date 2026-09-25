# =============================================================================# Stage 1: Base dependencies (shared)
# =============================================================================
FROM golang:1.25-alpine AS base

RUN apk add --no-cache git ca-certificates tzdata curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install golang-migrate CLI (dikunci versi untuk reproducibility)
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# =============================================================================
# Stage 2: Development (Air hot-reload)
# =============================================================================
FROM base AS dev

# Install Air untuk hot-reload
RUN go install github.com/air-verse/air@v1.61.7

WORKDIR /app

# Copy seluruh source code (akan di-override oleh volume mount di dev)
COPY . .

EXPOSE 8080

ENTRYPOINT ["./scripts/entrypoint.sh"]

# =============================================================================
# Stage 3: Builder (compile binary untuk production)
# =============================================================================
FROM base AS builder

WORKDIR /app

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

# =============================================================================
# Stage 4: Production runtime (image minimal)
# =============================================================================
FROM alpine:3.20 AS prod

RUN apk add --no-cache ca-certificates tzdata curl

ENV TZ=Asia/Jakarta

# Non-root user untuk security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy hanya artifact yang dibutuhkan dari builder
COPY --from=builder /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

.
RUN chmod +x ./entrypoint.sh

# Run as non-root
USER appuser

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]

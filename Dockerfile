# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install golang-migrate CLI (versi dikunci untuk reproducibility)
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

# Copy source code
COPY . .

# Build the binary (CGO disabled untuk binary static yang lebih kecil)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

# =============================================================================
# Stage 2: Runtime
# =============================================================================
FROM alpine:3.20

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Set timezone
ENV TZ=Asia/Jakarta

# Create non-root user untuk security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/server .

# Copy migrate CLI dari builder
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy entrypoint script
COPY --from=builder /app/scripts/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

# Run as non-root
USER appuser

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]

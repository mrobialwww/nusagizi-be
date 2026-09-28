#!/bin/sh
set -e

echo "============================================"
echo " NusaGizi Backend — Entrypoint"
echo " Environment: ${APP_ENV:-production}"
echo "============================================"

# Validate required environment variable
if [ -z "$DATABASE_URL" ]; then
  echo "[ERROR] DATABASE_URL is not set. Exiting."
  exit 1
fi

# Wait for PostgreSQL to be ready (secondary check after healthcheck)
# pg_isready does not support full URLs — parse the components manually
echo "[1/3] Waiting for database to be ready..."

# Extract host and port from DATABASE_URL (format: postgres://user:pass@host:port/db)
DB_HOST=$(echo "$DATABASE_URL" | sed -E 's|.*@([^:/]+).*|\1|')
DB_PORT=$(echo "$DATABASE_URL" | sed -E 's|.*@[^:]+:([0-9]+).*|\1|')
DB_PORT=${DB_PORT:-5432}

MAX_RETRIES=30
RETRY=0
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -q 2>/dev/null || [ $RETRY -ge $MAX_RETRIES ]; do
  RETRY=$((RETRY + 1))
  echo "  Attempt $RETRY/$MAX_RETRIES — retrying in 2s..."
  sleep 2
done

if [ $RETRY -ge $MAX_RETRIES ]; then
  echo "[ERROR] Database did not become ready in time. Exiting."
  exit 1
fi
echo "  Database is ready."

# Run database migrations
echo "[2/3] Running database migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up
echo "  Migrations applied successfully."

# Start server based on APP_ENV
echo "[3/3] Starting server..."

if [ "${APP_ENV}" = "development" ]; then
  echo "  Mode: DEVELOPMENT (Air hot-reload)"
  exec air -c .air.toml
else
  echo "  Mode: PRODUCTION (static binary)"
  exec ./server
fi

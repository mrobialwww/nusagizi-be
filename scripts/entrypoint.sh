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
echo "[1/3] Waiting for database to be ready..."
MAX_RETRIES=30
RETRY=0
until pg_isready -d "$DATABASE_URL" -q 2>/dev/null || [ $RETRY -ge $MAX_RETRIES ]; do
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

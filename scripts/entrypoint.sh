#!/bin/sh
set -e

echo "============================================"
echo " NusaGizi Backend — Entrypoint"
echo " Environment: ${APP_ENV:-production}"
echo "============================================"

# -----------------------------------------------------------------------------
# Validasi variabel wajib
# -----------------------------------------------------------------------------
if [ -z "$DATABASE_URL" ]; then
  echo "[ERROR] DATABASE_URL is not set. Exiting."
  exit 1
fi

# -----------------------------------------------------------------------------
# Tunggu hingga PostgreSQL siap (fallback jika healthcheck belum cukup)
# -----------------------------------------------------------------------------
echo "[1/3] Waiting for database to be ready..."
MAX_RETRIES=30
RETRY=0
until pg_isready -d "$DATABASE_URL" -q 2>/dev/null || [ $RETRY -ge $MAX_RETRIES ]; do
  RETRY=$((RETRY + 1))
  echo "  Attempt $RETRY/$MAX_RETRIES — DB not ready yet, retrying in 2s..."
  sleep 2
done

if [ $RETRY -ge $MAX_RETRIES ]; then
  echo "[ERROR] Database did not become ready in time. Exiting."
  exit 1
fi
echo "  Database is ready."

# -----------------------------------------------------------------------------
# Jalankan migrasi database
# -----------------------------------------------------------------------------
echo "[2/3] Running database migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up
echo "  Migrations applied successfully."

# -----------------------------------------------------------------------------
# Start server sesuai APP_ENV
# -----------------------------------------------------------------------------
echo "[3/3] Starting server..."

if [ "${APP_ENV}" = "development" ]; then
  echo "  Mode: DEVELOPMENT (Air hot-reload)"
  # Air membaca konfigurasi dari .air.toml di root project
  exec air -c .air.toml
else
  echo "  Mode: PRODUCTION (static binary)"
  exec ./server
fi

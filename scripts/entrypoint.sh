#!/bin/sh
set -e

echo "=== NusaGizi Backend Entrypoint ==="

# -----------------------------------------------------------------------------
# Jalankan migrasi database menggunakan golang-migrate
# Migrate diinstall saat build (lihat Dockerfile)
# -----------------------------------------------------------------------------
echo "[1/2] Running database migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up
echo "Migrations applied successfully."

# -----------------------------------------------------------------------------
# Start server
# -----------------------------------------------------------------------------
echo "[2/2] Starting backend server on port ${PORT:-8080}..."
exec ./server

#!/bin/bash
# =============================================================================
# scripts/init-db.sh
#
# Dijalankan otomatis oleh PostgreSQL saat container pertama kali dibuat
# (via /docker-entrypoint-initdb.d/).
#
# Tugasnya: buat database development jika belum ada.
# Database production sudah dibuat oleh POSTGRES_DB di env.
# =============================================================================
set -e

DB_DEV="${POSTGRES_DB_DEV:-nusagizi_dev}"
DB_USER="${POSTGRES_USER:-nusagizi}"

echo "[init-db] Checking if database '$DB_DEV' exists..."

# Buat DB dev jika belum ada
psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname "postgres" <<-EOSQL
    SELECT 'CREATE DATABASE "$DB_DEV"'
    WHERE NOT EXISTS (
        SELECT FROM pg_database WHERE datname = '$DB_DEV'
    )\gexec
EOSQL

echo "[init-db] Database '$DB_DEV' is ready."

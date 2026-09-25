#!/bin/bash
# =============================================================================
# scripts/init-db.sh
#
# Runs automatically on first PostgreSQL container startup
# via /docker-entrypoint-initdb.d/.
#
# Creates the dev database if it does not exist.
# The prod database is already created by POSTGRES_DB in the environment.
# =============================================================================
set -e

DB_DEV="${POSTGRES_DB_DEV:-nusagizi_dev}"
DB_USER="${POSTGRES_USER:-nusagizi}"

echo "[init-db] Checking if database '$DB_DEV' exists..."

# Create dev database only if it does not already exist (idempotent)
psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname "postgres" <<-EOSQL
    SELECT 'CREATE DATABASE "$DB_DEV"'
    WHERE NOT EXISTS (
        SELECT FROM pg_database WHERE datname = '$DB_DEV'
    )\gexec
EOSQL

echo "[init-db] Database '$DB_DEV' is ready."

#!/bin/bash
# =============================================================================
# scripts/init-db.sh
#
# Runs automatically on first PostgreSQL container startup
# via /docker-entrypoint-initdb.d/.
#
# Creates the staging database if it does not exist.
# The prod database is already created by POSTGRES_DB in the environment.
# =============================================================================
set -e

DB_STAGING="${POSTGRES_DB_STAGING:-nusagizi_staging}"
DB_USER="${POSTGRES_USER:-nusagizi}"

echo "[init-db] Checking if database '$DB_STAGING' exists..."

# Create staging database only if it does not already exist (idempotent)
psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname "postgres" <<-EOSQL
    SELECT 'CREATE DATABASE "$DB_STAGING"'
    WHERE NOT EXISTS (
        SELECT FROM pg_database WHERE datname = '$DB_STAGING'
    )\gexec
EOSQL

echo "[init-db] Database '$DB_STAGING' is ready."

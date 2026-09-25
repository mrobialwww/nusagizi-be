# =============================================================================
# Makefile — NusaGizi Backend
#
# Shortcut untuk operasi Docker dev/prod di VPS.
# =============================================================================

# Compose file combinations
COMPOSE_BASE   := docker compose -f docker-compose.yml
COMPOSE_PROD   := $(COMPOSE_BASE) -f docker-compose.prod.yml
COMPOSE_DEV    := $(COMPOSE_BASE) -f docker-compose.dev.yml

.PHONY: help \
        prod-up prod-down prod-logs prod-ps prod-restart prod-build \
        dev-up dev-down dev-logs dev-ps dev-restart dev-build \
        ps logs-postgres logs-minio \
        db-shell db-migrate db-migrate-dev \
        infra-up infra-down \
        clean prune

# Default target
help: ## Tampilkan daftar perintah yang tersedia
	@echo ""
	@echo "NusaGizi Backend — Makefile Commands"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Production
# =============================================================================

prod-build: ## Build image production
	$(COMPOSE_PROD) build backend-prod

prod-up: ## Jalankan semua service (production)
	$(COMPOSE_PROD) up -d

prod-down: ## Hentikan semua service (production)
	$(COMPOSE_PROD) down

prod-restart: ## Restart backend production saja
	$(COMPOSE_PROD) restart backend-prod

prod-logs: ## Lihat log backend production (follow)
	$(COMPOSE_PROD) logs -f backend-prod

prod-ps: ## Status container production
	$(COMPOSE_PROD) ps

# =============================================================================
# Development
# =============================================================================

dev-build: ## Build image development
	$(COMPOSE_DEV) build backend-dev

dev-up: ## Jalankan semua service (development)
	$(COMPOSE_DEV) up -d

dev-down: ## Hentikan semua service (development)
	$(COMPOSE_DEV) down

dev-restart: ## Restart backend development saja
	$(COMPOSE_DEV) restart backend-dev

dev-logs: ## Lihat log backend development (follow)
	$(COMPOSE_DEV) logs -f backend-dev

dev-ps: ## Status container development
	$(COMPOSE_DEV) ps

# =============================================================================
# Shared infrastructure (postgres + minio saja)
# =============================================================================

infra-up: ## Jalankan hanya postgres dan minio
	$(COMPOSE_BASE) up -d postgres minio

infra-down: ## Hentikan postgres dan minio
	$(COMPOSE_BASE) down postgres minio

# =============================================================================
# Monitoring
# =============================================================================

ps: ## Status semua container NusaGizi
	docker ps --filter "name=nusagizi"

logs-postgres: ## Lihat log PostgreSQL (follow)
	$(COMPOSE_BASE) logs -f postgres

logs-minio: ## Lihat log MinIO (follow)
	$(COMPOSE_BASE) logs -f minio

# =============================================================================
# Database utilities
# =============================================================================

db-shell: ## Masuk ke psql database production
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB:-nusagizi_prod}

db-shell-dev: ## Masuk ke psql database development
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB_DEV:-nusagizi_dev}

db-migrate: ## Jalankan migrasi di production (manual)
	$(COMPOSE_PROD) exec backend-prod migrate -path ./migrations -database "$$DATABASE_URL" up

db-migrate-dev: ## Jalankan migrasi di development (manual)
	$(COMPOSE_DEV) exec backend-dev migrate -path ./migrations -database "$$DATABASE_URL" up

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Hentikan semua container dan hapus anonymous volumes
	$(COMPOSE_PROD) down -v --remove-orphans 2>/dev/null || true
	$(COMPOSE_DEV) down -v --remove-orphans 2>/dev/null || true

prune: ## Hapus image yang tidak terpakai (jalankan dengan hati-hati)
	docker image prune -f
	docker volume prune -f

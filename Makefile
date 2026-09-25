# =============================================================================
# Makefile — NusaGizi Backend
# Shortcuts for Docker dev/prod operations on the VPS.
# =============================================================================

# Compose file combinations
COMPOSE_BASE   := docker compose -f docker-compose.yml
COMPOSE_PROD   := $(COMPOSE_BASE) -f docker-compose.prod.yml
COMPOSE_DEV    := $(COMPOSE_BASE) -f docker-compose.dev.yml

.PHONY: help \
        prod-up prod-down prod-logs prod-ps prod-restart prod-build \
        dev-up dev-down dev-logs dev-ps dev-restart dev-build \
        ps logs-postgres logs-minio \
        db-shell db-shell-dev db-migrate db-migrate-dev \
        infra-up infra-down \
        clean prune

# Default target
help: ## Show available commands
	@echo ""
	@echo "NusaGizi Backend — Makefile Commands"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Production
# =============================================================================

prod-build: ## Build production image
	$(COMPOSE_PROD) build backend-prod

prod-up: ## Start all production services
	$(COMPOSE_PROD) up -d

prod-down: ## Stop all production services
	$(COMPOSE_PROD) down

prod-restart: ## Restart production backend only
	$(COMPOSE_PROD) restart backend-prod

prod-logs: ## Follow production backend logs
	$(COMPOSE_PROD) logs -f backend-prod

prod-ps: ## Show production container status
	$(COMPOSE_PROD) ps

# =============================================================================
# Development
# =============================================================================

dev-build: ## Build development image
	$(COMPOSE_DEV) build backend-dev

dev-up: ## Start all development services
	$(COMPOSE_DEV) up -d

dev-down: ## Stop all development services
	$(COMPOSE_DEV) down

dev-restart: ## Restart development backend only
	$(COMPOSE_DEV) restart backend-dev

dev-logs: ## Follow development backend logs
	$(COMPOSE_DEV) logs -f backend-dev

dev-ps: ## Show development container status
	$(COMPOSE_DEV) ps

# =============================================================================
# Shared infrastructure (postgres + minio only)
# =============================================================================

infra-up: ## Start postgres and minio only
	$(COMPOSE_BASE) up -d postgres minio

infra-down: ## Stop postgres and minio
	$(COMPOSE_BASE) down postgres minio

# =============================================================================
# Monitoring
# =============================================================================

ps: ## Show all NusaGizi container status
	docker ps --filter "name=nusagizi"

logs-postgres: ## Follow PostgreSQL logs
	$(COMPOSE_BASE) logs -f postgres

logs-minio: ## Follow MinIO logs
	$(COMPOSE_BASE) logs -f minio

# =============================================================================
# Database utilities
# =============================================================================

db-shell: ## Open psql shell for production database
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB:-nusagizi_prod}

db-shell-dev: ## Open psql shell for development database
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB_DEV:-nusagizi_dev}

db-migrate: ## Run migrations on production (manual)
	$(COMPOSE_PROD) exec backend-prod migrate -path ./migrations -database "$$DATABASE_URL" up

db-migrate-dev: ## Run migrations on development (manual)
	$(COMPOSE_DEV) exec backend-dev migrate -path ./migrations -database "$$DATABASE_URL" up

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Stop all containers and remove anonymous volumes
	$(COMPOSE_PROD) down -v --remove-orphans 2>/dev/null || true
	$(COMPOSE_DEV) down -v --remove-orphans 2>/dev/null || true

prune: ## Remove unused images and volumes (use with caution)
	docker image prune -f
	docker volume prune -f

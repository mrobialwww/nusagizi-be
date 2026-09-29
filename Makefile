# =============================================================================
# Makefile — NusaGizi Backend
# Shortcuts for Docker staging/prod operations on the VPS.
# =============================================================================

# Auto-detect which base file to use based on current directory name.
# - nusagizi-be-prod   : uses docker-compose.yml (owns postgres)
# - nusagizi-be-staging: uses docker-compose.base-staging.yml (postgres is external)
CURRENT_DIR := $(notdir $(CURDIR))
ifeq ($(CURRENT_DIR), nusagizi-be-staging)
  COMPOSE_BASE := docker compose -f docker-compose.base-staging.yml
else
  COMPOSE_BASE := docker compose -f docker-compose.yml
endif

COMPOSE_PROD    := $(COMPOSE_BASE) -f docker-compose.prod.yml
COMPOSE_STAGING := $(COMPOSE_BASE) -f docker-compose.staging.yml

.PHONY: help \
        prod-up prod-down prod-logs prod-ps prod-restart prod-build \
        staging-up staging-down staging-logs staging-ps staging-restart staging-build \
        ps logs-postgres \
        db-shell db-shell-staging db-migrate db-migrate-staging \
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
# Staging
# =============================================================================

staging-build: ## Build staging image
	$(COMPOSE_STAGING) build backend-staging

staging-up: ## Start all staging services
	$(COMPOSE_STAGING) up -d

staging-down: ## Stop all staging services
	$(COMPOSE_STAGING) down

staging-restart: ## Restart staging backend only
	$(COMPOSE_STAGING) restart backend-staging

staging-logs: ## Follow staging backend logs
	$(COMPOSE_STAGING) logs -f backend-staging

staging-ps: ## Show staging container status
	$(COMPOSE_STAGING) ps

# =============================================================================
# Shared infrastructure (postgres only)
# =============================================================================

infra-up: ## Start postgres (run from nusagizi-be-prod only)
	docker compose -f docker-compose.yml up -d postgres

infra-down: ## Stop postgres (run from nusagizi-be-prod only)
	docker compose -f docker-compose.yml down postgres

# =============================================================================
# Monitoring
# =============================================================================

ps: ## Show all NusaGizi container status
	docker ps --filter "name=nusagizi"

logs-postgres: ## Follow PostgreSQL logs
	$(COMPOSE_BASE) logs -f postgres

# =============================================================================
# Database utilities
# =============================================================================

db-shell: ## Open psql shell for production database
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB:-nusagizi_prod}

db-shell-staging: ## Open psql shell for staging database
	docker exec -it nusagizi-postgres psql -U $${POSTGRES_USER:-nusagizi} -d $${POSTGRES_DB_STAGING:-nusagizi_staging}

db-migrate: ## Run migrations on production (manual)
	$(COMPOSE_PROD) exec backend-prod migrate -path ./migrations -database "$$DATABASE_URL" up

db-migrate-staging: ## Run migrations on staging (manual)
	$(COMPOSE_STAGING) exec backend-staging migrate -path ./migrations -database "$$DATABASE_URL" up

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Stop all containers and remove anonymous volumes
	$(COMPOSE_PROD) down -v --remove-orphans 2>/dev/null || true
	$(COMPOSE_STAGING) down -v --remove-orphans 2>/dev/null || true

prune: ## Remove unused images and volumes (use with caution)
	docker image prune -f
	docker volume prune -f

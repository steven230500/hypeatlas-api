# ========= Variables =========
POSTGRES_URL ?= postgres://user:pass@localhost:5432/hypeatlas?sslmode=disable

# Docker Compose (ajusta si usas "docker compose")
DC := docker-compose

# Swagger (swaggo) — usamos go run pinneado para evitar desalineación de versiones
SWAG_VERSION ?= v1.16.3
SWAG_CMD     := go run github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
SWAG_MAIN    ?= cmd/api/main.go
SWAG_OUT     ?= docs

# Goose (migraciones) — asumimos que ya lo tienes instalado
GOOSE := $(shell go env GOPATH)/bin/goose
MIGR_DIR := ./db/migrations

# ========= Ayuda =========
.PHONY: help
help:
	@echo "Targets:"
	@echo "  db-up            - Levanta servicios con docker-compose"
	@echo "  db-down          - Baja servicios con docker-compose"
	@echo "  migrate-up       - Aplica migraciones (goose up)"
	@echo "  migrate-down     - Revierte última migración (goose down)"
	@echo "  migrate-status   - Estado de migraciones"
	@echo "  docs             - Genera OpenAPI (swag init) en ./docs"
	@echo "  docs-clean       - Elimina ./docs generado"
	@echo "  docs-print       - Imprime versión de swag y valida doc.json"
	@echo "  run              - Genera docs y ejecuta API (STORAGE=postgres)"
	@echo "  worker           - Ejecuta worker"

# ========= Docker =========
.PHONY: db-up db-down
db-up:
	$(DC) up -d

db-down:
	$(DC) down

# ========= Migraciones =========
.PHONY: migrate-up migrate-down migrate-status
migrate-up:
	$(GOOSE) -dir $(MIGR_DIR) postgres "$(POSTGRES_URL)" up

migrate-down:
	$(GOOSE) -dir $(MIGR_DIR) postgres "$(POSTGRES_URL)" down

migrate-status:
	$(GOOSE) -dir $(MIGR_DIR) postgres "$(POSTGRES_URL)" status

# ========= Swagger / OpenAPI =========
.PHONY: docs docs-clean docs-print
docs:
	rm -rf $(SWAG_OUT)
	$(SWAG_CMD) init --parseDependency --parseInternal -g $(SWAG_MAIN) -o $(SWAG_OUT)

docs-clean:
	rm -rf $(SWAG_OUT)

docs-print:
	@echo "swag version (requested): $(SWAG_VERSION)"
	@$(SWAG_CMD) -v
	@test -f $(SWAG_OUT)/doc.json && echo "docs/doc.json OK" || (echo "docs/doc.json NO EXISTE"; exit 1)

# ========= App =========
.PHONY: run worker
run: docs
	STORAGE=postgres POSTGRES_URL=$(POSTGRES_URL) go run ./cmd/api

worker:
	POSTGRES_URL=$(POSTGRES_URL) go run ./cmd/worker

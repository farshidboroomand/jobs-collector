include .env
export

GREEN  := \033[32m
YELLOW := \033[33m
CYAN   := \033[36m
RESET  := \033[0m

.DEFAULT_GOAL := help

COMPOSE_FILE := deployments/docker/docker-compose.yml
DOCKER_COMPOSE := docker compose --env-file .env -f $(COMPOSE_FILE)

MIGRATE_DB := $(DB_CONNECTION)://$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)?multiStatements=true&parseTime=true
VERSION ?= 1

.PHONY: help build lint lint-fix up down restart logs reset-db migrate-create migrate-up migrate-down migrate-force

help:
	@echo ""
	@echo "$(CYAN)Available commands$(RESET)"
	@echo ""
	@echo "  make up               Start all services"
	@echo "  make down             Stop all services"
	@echo "  make restart          Restart services"
	@echo "  make logs             Follow container logs"
	@echo "  make reset-db         Reset database (destroy volume)"
	@echo ""
	@echo "  make build            Build Go binaries"
	@echo "  make lint             Run golangci-lint"
	@echo "  make lint-fix         Auto-format code"
	@echo ""
	@echo "  make migrate-create name=create_jobs_table"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make migrate-force VERSION=1"
	@echo ""

build:
	@echo "$(GREEN)Building Go binaries...$(RESET)"
	@go build -o bin/ ./cmd/...
	@echo "$(GREEN)Build completed$(RESET)"

lint:
	@echo "$(GREEN)Running golangci-lint...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)Linting complete$(RESET)"

lint-fix:
	@echo "$(GREEN)Formatting code...$(RESET)"
	@golangci-lint fmt
	@echo "$(GREEN)Formatting complete$(RESET)"

up:
	@echo "$(CYAN)Starting services...$(RESET)"
	@$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Services are running$(RESET)"

down:
	@echo "$(YELLOW)Stopping services...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(YELLOW)Services stopped$(RESET)"

restart:
	@echo "$(CYAN)Restarting services...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Services restarted$(RESET)"

logs:
	@$(DOCKER_COMPOSE) logs -f

reset-db:
	@echo "$(YELLOW)Resetting database...$(RESET)"
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Database reset completed$(RESET)"

migrate-create:
	@migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	@echo "$(CYAN)Running migrations up...$(RESET)"
	@$(DOCKER_COMPOSE) run --rm migrate \
	-path=/migrations \
	-database "$(MIGRATE_DB)" \
	up

migrate-down:
	@echo "$(YELLOW)Rolling back migration...$(RESET)"
	@$(DOCKER_COMPOSE) run --rm migrate \
	-path=/migrations \
	-database "$(MIGRATE_DB)" \
	down 1

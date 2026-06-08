GREEN  := \033[32m
YELLOW := \033[33m
CYAN   := \033[36m
RESET  := \033[0m

.DEFAULT_GOAL := help

COMPOSE_FILE := deployments/docker/docker-compose.yml

.PHONY: build lint lint-fix up down logs

build:
	@echo "$(GREEN)Building the project...$(RESET)"
	@go build -o bin/ ./cmd/...
	@echo "$(GREEN)Build completed$(RESET)"

lint:
	@echo "$(GREEN)Linting code with golangci-lint...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)Linting complete$(RESET)"

lint-fix:
	@echo "$(GREEN)Formatting code with golangci-lint...$(RESET)"
	@golangci-lint fmt
	@echo "$(GREEN)Code formatting complete$(RESET)"

up:
	@echo "$(CYAN)Starting services...$(RESET)"
	@docker compose --env-file .env -f $(COMPOSE_FILE) up -d --build
	@echo "$(GREEN)Services are up$(RESET)"

down:
	@echo "$(YELLOW)Stopping services...$(RESET)"
	@docker compose -f $(COMPOSE_FILE) down

logs:
	@docker compose -f $(COMPOSE_FILE) logs -f

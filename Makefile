GREEN  := \033[32m
YELLOW := \033[33m
CYAN   := \033[36m
RESET  := \033[0m

.DEFAULT_GOAL := help

.PHONY: build lint lint-fix

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
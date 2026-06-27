# Makefile untuk ESQL-AST Tool - Versi Shell Script

.PHONY: help build clean test baseline diff update list

BINARY = esql-ast
GREEN = \033[0;32m
RED = \033[0;31m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m

.DEFAULT_GOAL := help

help:
	@echo "$(BLUE)ESQL-AST Tool - Commands$(NC)"
	@echo ""
	@echo "  $(GREEN)make build$(NC)       - Build binary"
	@echo "  $(GREEN)make baseline$(NC)    - Generate baseline for all files"
	@echo "  $(GREEN)make test$(NC)        - Test all files against baseline"
	@echo "  $(GREEN)make diff$(NC)        - Show differences"
	@echo "  $(GREEN)make update$(NC)      - Update baseline (with confirmation)"
	@echo "  $(GREEN)make clean$(NC)       - Clean build and outputs"
	@echo "  $(GREEN)make list$(NC)        - List all test files"

build:
	@echo "$(YELLOW)Building...$(NC)"
	@go build -o $(BINARY) cmd/esql-ast/main.go
	@echo "$(GREEN)✅ Build done$(NC)"

clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -f $(BINARY)
	@rm -rf tests/output
	@rm -rf tests/baseline
	@rm -rf tests/diff
	@echo "$(GREEN)✅ Clean done$(NC)"

list:
	@echo "$(BLUE)Test files:$(NC)"
	@for f in examples/*.esql; do echo "  - $$(basename $$f .esql)"; done

baseline: build
	@chmod +x scripts/baseline.sh
	@./scripts/baseline.sh

test: build
	@chmod +x scripts/test.sh
	@./scripts/test.sh

diff: build
	@chmod +x scripts/diff.sh
	@./scripts/diff.sh

update: build
	@echo "$(RED)⚠️  WARNING: This will OVERWRITE all baselines!$(NC)"
	@printf "Are you sure? (y/N) "
	@read ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		./scripts/baseline.sh; \
		echo "$(GREEN)✅ Baselines updated!$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled$(NC)"; \
	fi
# Makefile untuk ESQL-AST Tool

.PHONY: help build clean test test-all test-sample test-cast test-case test-regression baseline update diff validate coverage

# Variables
BINARY = esql-ast
BASELINE_DIR = tests/baseline
OUTPUT_DIR = tests/output
DIFF_DIR = tests/diff

# Colors
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m

# ============================================
# HELP
# ============================================

help:
	@echo "$(BLUE)ESQL-AST Tool - Makefile Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Build:$(NC)"
	@echo "  make build         - Build esql-ast binary"
	@echo "  make clean         - Remove binary and test outputs"
	@echo ""
	@echo "$(YELLOW)Test:$(NC)"
	@echo "  make test          - Run all tests"
	@echo "  make test-sample   - Test sample.esql only"
	@echo "  make test-cast     - Test CAST features"
	@echo "  make test-case     - Test CASE features"
	@echo "  make test-all      - Run all tests with diff"
	@echo "  make test-quick    - Quick test without diff"
	@echo "  make test-regression - Regression test (sample only)"
	@echo ""
	@echo "$(YELLOW)Baseline:$(NC)"
	@echo "  make baseline      - Generate all baselines"
	@echo "  make baseline-sample - Generate sample baseline only"
	@echo "  make baseline-cast - Generate CAST baseline"
	@echo "  make baseline-case - Generate CASE baseline"
	@echo ""
	@echo "$(YELLOW)Diff:$(NC)"
	@echo "  make diff          - Show diff between baseline and current"
	@echo "  make diff-sample   - Show diff for sample.esql only"
	@echo ""
	@echo "$(YELLOW)Utils:$(NC)"
	@echo "  make validate      - Validate all ESQL files"
	@echo "  make coverage      - Generate test coverage"
	@echo "  make update        - Update baseline (after verification)"
	@echo "  make dirs          - Create test directories"

# ============================================
# BUILD
# ============================================

build:
	@echo "$(YELLOW)Building $(BINARY)...$(NC)"
	go build -o $(BINARY) cmd/esql-ast/main.go
	@echo "$(GREEN)Build successful!$(NC)"

clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -f $(BINARY)
	rm -rf $(OUTPUT_DIR)
	rm -rf $(DIFF_DIR)
	rm -f /tmp/esql-*.txt
	@echo "$(GREEN)Clean complete!$(NC)"

dirs:
	@mkdir -p $(BASELINE_DIR)
	@mkdir -p $(OUTPUT_DIR)
	@mkdir -p $(DIFF_DIR)

# ============================================
# BASELINE GENERATION
# ============================================

baseline: build dirs baseline-sample baseline-cast baseline-case
	@echo "$(GREEN)All baselines generated!$(NC)"

baseline-sample: build dirs
	@echo "$(YELLOW)Generating baseline for sample.esql...$(NC)"
	@echo "  - Pretty..."
	./$(BINARY) -f examples/sample.esql -pretty > $(BASELINE_DIR)/sample.pretty.txt 2>&1 || true
	@echo "  - JSON..."
	./$(BINARY) -f examples/sample.esql -json > $(BASELINE_DIR)/sample.json.txt 2>&1 || true
	@echo "  - Generate..."
	./$(BINARY) -f examples/sample.esql -generate > $(BASELINE_DIR)/sample.generate.txt 2>&1 || true
	@echo "  - Analyze..."
	./$(BINARY) -f examples/sample.esql -analyze > $(BASELINE_DIR)/sample.analyze.txt 2>&1 || true
	@echo "$(GREEN)Sample baseline done!$(NC)"

baseline-cast: build dirs
	@echo "$(YELLOW)Generating baseline for CAST tests...$(NC)"
	@echo "  - test_cast.esql pretty..."
	./$(BINARY) -f examples/test_cast.esql -pretty > $(BASELINE_DIR)/cast.pretty.txt 2>&1 || true
	@echo "  - test_cast.esql json..."
	./$(BINARY) -f examples/test_cast.esql -json > $(BASELINE_DIR)/cast.json.txt 2>&1 || true
	@echo "  - test_cast.esql generate..."
	./$(BINARY) -f examples/test_cast.esql -generate > $(BASELINE_DIR)/cast.generate.txt 2>&1 || true
	@echo "  - test_cast.esql analyze..."
	./$(BINARY) -f examples/test_cast.esql -analyze > $(BASELINE_DIR)/cast.analyze.txt 2>&1 || true
	@echo "  - test_nested_cast.esql pretty..."
	./$(BINARY) -f examples/test_nested_cast.esql -pretty > $(BASELINE_DIR)/nested_cast.pretty.txt 2>&1 || true
	@echo "  - test_nested_cast.esql generate..."
	./$(BINARY) -f examples/test_nested_cast.esql -generate > $(BASELINE_DIR)/nested_cast.generate.txt 2>&1 || true
	@echo "$(GREEN)CAST baseline done!$(NC)"

baseline-case: build dirs
	@echo "$(YELLOW)Generating baseline for CASE tests...$(NC)"
	@echo "  - test_case_simple_only.esql..."
	./$(BINARY) -f examples/test_case_simple_only.esql -pretty > $(BASELINE_DIR)/case_simple.pretty.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_simple_only.esql -generate > $(BASELINE_DIR)/case_simple.generate.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_simple_only.esql -analyze > $(BASELINE_DIR)/case_simple.analyze.txt 2>&1 || true
	@echo "  - test_case_searched_only.esql..."
	./$(BINARY) -f examples/test_case_searched_only.esql -pretty > $(BASELINE_DIR)/case_searched.pretty.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_searched_only.esql -generate > $(BASELINE_DIR)/case_searched.generate.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_searched_only.esql -analyze > $(BASELINE_DIR)/case_searched.analyze.txt 2>&1 || true
	@echo "  - test_case_nested_if.esql..."
	./$(BINARY) -f examples/test_case_nested_if.esql -pretty > $(BASELINE_DIR)/case_nested_if.pretty.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_nested_if.esql -generate > $(BASELINE_DIR)/case_nested_if.generate.txt 2>&1 || true
	./$(BINARY) -f examples/test_case_nested_if.esql -analyze > $(BASELINE_DIR)/case_nested_if.analyze.txt 2>&1 || true
	@echo "  - test_case.esql..."
	./$(BINARY) -f examples/test_case.esql -pretty > $(BASELINE_DIR)/case_full.pretty.txt 2>&1 || true
	./$(BINARY) -f examples/test_case.esql -generate > $(BASELINE_DIR)/case_full.generate.txt 2>&1 || true
	./$(BINARY) -f examples/test_case.esql -analyze > $(BASELINE_DIR)/case_full.analyze.txt 2>&1 || true
	@echo "$(GREEN)CASE baseline done!$(NC)"

# ============================================
# RUN TESTS
# ============================================

# Helper function untuk run test
define run_test
	@echo "$(YELLOW)Testing: $1 - $2$(NC)"
	@mkdir -p $(dir $(OUTPUT_DIR)/$3)
	./$(BINARY) -f $1 $2 > $(OUTPUT_DIR)/$3 2>&1 || true
	@if [ -f "$(BASELINE_DIR)/$3" ]; then \
		if diff -q $(BASELINE_DIR)/$3 $(OUTPUT_DIR)/$3 > /dev/null 2>&1; then \
			echo "  $(GREEN)✅ PASSED: $3$(NC)"; \
		else \
			echo "  $(RED)❌ FAILED: $3$(NC)"; \
			echo "  Diff saved to $(DIFF_DIR)/$3.diff"; \
			mkdir -p $(DIFF_DIR); \
			diff -u $(BASELINE_DIR)/$3 $(OUTPUT_DIR)/$3 > $(DIFF_DIR)/$3.diff 2>&1 || true; \
			FAILED=1; \
		fi \
	else \
		echo "  $(YELLOW)⚠️  No baseline for $3, creating...$(NC)"; \
		cp $(OUTPUT_DIR)/$3 $(BASELINE_DIR)/$3; \
	fi
endef

test: build dirs test-sample test-cast test-case
	@echo ""
	@echo "$(GREEN)All tests completed!$(NC)"

test-regression: build dirs
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)REGRESSION TEST - sample.esql only$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	$(call run_test,examples/sample.esql,-pretty,sample.pretty.txt)
	$(call run_test,examples/sample.esql,-json,sample.json.txt)
	$(call run_test,examples/sample.esql,-generate,sample.generate.txt)
	$(call run_test,examples/sample.esql,-analyze,sample.analyze.txt)

test-sample: build dirs
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)Testing SAMPLE.ESQL$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	$(call run_test,examples/sample.esql,-pretty,sample.pretty.txt)
	$(call run_test,examples/sample.esql,-json,sample.json.txt)
	$(call run_test,examples/sample.esql,-generate,sample.generate.txt)
	$(call run_test,examples/sample.esql,-analyze,sample.analyze.txt)

test-cast: build dirs
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)Testing CAST Features$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	$(call run_test,examples/test_cast.esql,-pretty,cast.pretty.txt)
	$(call run_test,examples/test_cast.esql,-json,cast.json.txt)
	$(call run_test,examples/test_cast.esql,-generate,cast.generate.txt)
	$(call run_test,examples/test_cast.esql,-analyze,cast.analyze.txt)
	$(call run_test,examples/test_nested_cast.esql,-pretty,nested_cast.pretty.txt)
	$(call run_test,examples/test_nested_cast.esql,-generate,nested_cast.generate.txt)

test-case: build dirs
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)Testing CASE Features$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	$(call run_test,examples/test_case_simple_only.esql,-pretty,case_simple.pretty.txt)
	$(call run_test,examples/test_case_simple_only.esql,-generate,case_simple.generate.txt)
	$(call run_test,examples/test_case_simple_only.esql,-analyze,case_simple.analyze.txt)
	$(call run_test,examples/test_case_searched_only.esql,-pretty,case_searched.pretty.txt)
	$(call run_test,examples/test_case_searched_only.esql,-generate,case_searched.generate.txt)
	$(call run_test,examples/test_case_searched_only.esql,-analyze,case_searched.analyze.txt)
	$(call run_test,examples/test_case_nested_if.esql,-pretty,case_nested_if.pretty.txt)
	$(call run_test,examples/test_case_nested_if.esql,-generate,case_nested_if.generate.txt)
	$(call run_test,examples/test_case_nested_if.esql,-analyze,case_nested_if.analyze.txt)
	$(call run_test,examples/test_case.esql,-pretty,case_full.pretty.txt)
	$(call run_test,examples/test_case.esql,-generate,case_full.generate.txt)
	$(call run_test,examples/test_case.esql,-analyze,case_full.analyze.txt)

test-all: build dirs test-sample test-cast test-case
	@echo ""
	@echo "$(BLUE)========================================$(NC)"
	@echo "$(BLUE)ALL TESTS COMPLETED$(NC)"
	@echo "$(BLUE)========================================$(NC)"
	@if [ -z "$$FAILED" ]; then \
		echo "$(GREEN)✅ All tests passed!$(NC)"; \
	else \
		echo "$(RED)❌ Some tests failed! Check $(DIFF_DIR)/*.diff$(NC)"; \
		exit 1; \
	fi

test-quick: build
	@echo "$(YELLOW)Quick test (no diff)...$(NC)"
	./$(BINARY) -f examples/sample.esql -pretty
	./$(BINARY) -f examples/test_cast.esql -pretty
	./$(BINARY) -f examples/test_case_simple_only.esql -pretty

# ============================================
# DIFF VIEW
# ============================================

diff: build dirs
	@echo "$(YELLOW)Generating diffs...$(NC)"
	@$(MAKE) diff-sample
	@$(MAKE) diff-cast
	@$(MAKE) diff-case

diff-sample: build dirs
	@echo "$(BLUE)Diff for sample.esql:$(NC)"
	@for f in sample.pretty sample.json sample.generate sample.analyze; do \
		if [ -f "$(BASELINE_DIR)/$$f.txt" ]; then \
			echo "  $$f:"; \
			./$(BINARY) -f examples/sample.esql -$$(echo $$f | cut -d. -f2) > $(OUTPUT_DIR)/$$f.txt 2>&1 || true; \
			diff -u $(BASELINE_DIR)/$$f.txt $(OUTPUT_DIR)/$$f.txt || echo "    No changes"; \
		fi; \
	done

diff-cast: build dirs
	@echo "$(BLUE)Diff for CAST tests:$(NC)"
	@for f in cast.pretty cast.json cast.generate cast.analyze nested_cast.pretty nested_cast.generate; do \
		if [ -f "$(BASELINE_DIR)/$$f.txt" ]; then \
			echo "  $$f:"; \
			./$(BINARY) -f examples/test_cast.esql -$$(echo $$f | cut -d. -f2) > $(OUTPUT_DIR)/$$f.txt 2>&1 || true; \
			diff -u $(BASELINE_DIR)/$$f.txt $(OUTPUT_DIR)/$$f.txt || echo "    No changes"; \
		fi; \
	done

diff-case: build dirs
	@echo "$(BLUE)Diff for CASE tests:$(NC)"
	@for f in case_simple.pretty case_simple.generate case_simple.analyze case_searched.pretty case_searched.generate case_searched.analyze case_nested_if.pretty case_nested_if.generate case_nested_if.analyze case_full.pretty case_full.generate case_full.analyze; do \
		if [ -f "$(BASELINE_DIR)/$$f.txt" ]; then \
			echo "  $$f:"; \
			./$(BINARY) -f examples/test_case_simple_only.esql -$$(echo $$f | cut -d. -f2) > $(OUTPUT_DIR)/$$f.txt 2>&1 || true; \
			diff -u $(BASELINE_DIR)/$$f.txt $(OUTPUT_DIR)/$$f.txt || echo "    No changes"; \
		fi; \
	done

# ============================================
# UPDATE BASELINE
# ============================================

update: build
	@echo "$(YELLOW)Updating baselines...$(NC)"
	@echo "$(RED)WARNING: This will overwrite existing baselines!$(NC)"
	@printf "Are you sure? (y/N) "
	@read answer; \
	if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
		$(MAKE) baseline; \
		echo "$(GREEN)Baselines updated!$(NC)"; \
	else \
		echo "$(YELLOW)Update cancelled.$(NC)"; \
	fi

# ============================================
# VALIDATE ESQL FILES
# ============================================

validate: build
	@echo "$(YELLOW)Validating all ESQL files...$(NC)"
	@for f in examples/*.esql; do \
		echo "Validating $$f..."; \
		./$(BINARY) -f $$f -validate || echo "  $(RED)❌ Validation failed for $$f$(NC)"; \
	done
	@echo "$(GREEN)Validation complete!$(NC)"

# ============================================
# COVERAGE
# ============================================

coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# ============================================
# CLEAN ALL (include baseline)
# ============================================

clean-all: clean
	@echo "$(YELLOW)Removing baselines...$(NC)"
	rm -rf $(BASELINE_DIR)
	@echo "$(GREEN)All cleaned!$(NC)"

# ============================================
# DEFAULT
# ============================================

.DEFAULT_GOAL := help
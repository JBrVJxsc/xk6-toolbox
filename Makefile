# Makefile for xk6-toolbox extension

# Variables
MODULE_NAME = github.com/xuzhang/xk6-toolbox
K6_VERSION = v1.0.0
XK6_VERSION = v0.20.1
BUILD_DIR = build
K6_BINARY = $(BUILD_DIR)/k6

# Colors for output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: all build clean test test-go test-k6 help

# Default target
all: build test

# Help target
help:
	@echo "$(GREEN)Available targets:$(NC)"
	@echo "  $(YELLOW)all$(NC)        - Build extension and run all tests"
	@echo "  $(YELLOW)build$(NC)      - Build k6 binary with toolbox extension"
	@echo "  $(YELLOW)test$(NC)       - Run all tests (Go + k6)"
	@echo "  $(YELLOW)test-go$(NC)    - Run Go unit tests only"
	@echo "  $(YELLOW)test-k6$(NC)    - Run k6 JavaScript tests only"
	@echo "  $(YELLOW)clean$(NC)      - Clean build artifacts"
	@echo "  $(YELLOW)help$(NC)       - Show this help message"

# Build k6 binary with extension
build:
	@echo "$(GREEN)Building k6 with toolbox extension...$(NC)"
	@mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && xk6 build $(K6_VERSION) --with $(MODULE_NAME)=../
	@echo "$(GREEN)✓ Build complete: $(K6_BINARY)$(NC)"

# Run all tests
test: test-go test-k6
	@echo "$(GREEN)✓ All tests completed successfully!$(NC)"

# Run Go unit tests
test-go:
	@echo "$(GREEN)Running Go unit tests...$(NC)"
	go test -v ./...
	@echo "$(GREEN)✓ Go tests completed$(NC)"

# Run k6 JavaScript tests
test-k6: build
	@echo "$(GREEN)Running k6 JavaScript tests...$(NC)"
	@if [ -f "toolbox_k6_test.js" ]; then \
		$(K6_BINARY) run toolbox_k6_test.js; \
	else \
		echo "$(RED)Error: toolbox_k6_test.js not found$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ k6 tests completed$(NC)"

# Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ Clean completed$(NC)" 
# Authorizer Go SDK - Makefile
#
# Prerequisites for integration tests:
# 1. Start the authorizer server: make docker-up
# 2. Run tests: make test

# Docker image for authorizer server
AUTHORIZER_IMAGE := lakhansamani/authorizer:latest
AUTHORIZER_CONTAINER := authorizer-test

.PHONY: docker-up docker-down test

# Start authorizer in Docker for integration testing
docker-up:
	@if docker ps -q -f name=^/$(AUTHORIZER_CONTAINER)$$ | grep -q .; then \
		echo "Authorizer container already running"; \
	else \
		echo "Starting authorizer container..."; \
		docker run -d --rm \
			--name $(AUTHORIZER_CONTAINER) \
			-p 8080:8080 \
			$(AUTHORIZER_IMAGE) \
			--database-type=sqlite \
			--database-url=test.db \
			--jwt-type=HS256 \
			--jwt-secret=test \
			--admin-secret=admin \
			--client-id=123456 \
			--client-secret=secret; \
		echo "Waiting for authorizer to be ready..."; \
		sleep 3; \
		echo "Authorizer is running at http://localhost:8080"; \
	fi

# Stop the authorizer container
docker-down:
	@echo "Stopping authorizer container..."
	@docker stop $(AUTHORIZER_CONTAINER) 2>/dev/null || true

# Run integration tests - starts authorizer if needed, runs tests, then stops container
test: docker-up
	@echo "Running integration tests..."
	@go test -v ./test/ -count=1; \
	EXIT_CODE=$$?; \
	$(MAKE) docker-down; \
	exit $$EXIT_CODE

# Run tests only (assumes authorizer is already running - run 'make docker-up' first)
test-only:
	@echo "Running integration tests..."
	@go test -v ./test/ -count=1

COMPOSE=docker-compose -f docker-compose.yaml

.PHONY: .setup-db .teardown-db

.PHONY: test test-integration

test:
	go test ./...

test-integration: .setup-db
	go test -tags=integration ./...
	@$(MAKE) .teardown-db

.setup-db:
	@echo "🚀 Starting PostgreSQL on port 5433..."
	@$(COMPOSE) up --wait -d

.teardown-db:
	@echo "🧹 Stopping PostgreSQL..."
	@$(COMPOSE) down --volumes --remove-orphans

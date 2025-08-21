COMPOSE=docker-compose -f docker-compose.yaml

.PHONY: .setup-db .teardown-db

.PHONY: test test-integration

test:
	go test -v -tags=unit ./...

test-integration: .setup-db
	go test -v -tags=integration ./...
	@bash -c 'set -e; trap "make .teardown-db" EXIT; go test -tags=integration ./...'

.setup-db:
	@echo "🚀 Starting PostgreSQL on port 5433..."
	@$(COMPOSE) up --wait -d

.teardown-db:
	@echo "🧹 Stopping PostgreSQL..."
	@$(COMPOSE) down --volumes --remove-orphans

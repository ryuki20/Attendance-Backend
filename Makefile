.PHONY: help up down build logs migrate-up migrate-down migrate-create clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start docker containers
	docker-compose up -d

down: ## Stop docker containers
	docker-compose down

build: ## Build docker containers
	docker-compose build

logs: ## Show docker logs
	docker-compose logs -f app

migrate-up: ## Run database migrations up
	docker-compose exec app migrate -path=/app/migrations -database "postgresql://attendance_user:attendance_password@db:5432/attendance_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	docker-compose exec app migrate -path=/app/migrations -database "postgresql://attendance_user:attendance_password@db:5432/attendance_db?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	docker-compose exec app migrate create -ext sql -dir /app/migrations -seq $(name)

clean: ## Clean up docker volumes and containers
	docker-compose down -v
	rm -rf tmp/

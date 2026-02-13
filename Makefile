
.DEFAULT_GOAL := help
.PHONY: help migrate migrate-up migrate-down migrate-force migrate-version seed dev infra-up infra-down gen-docs
MIGRATION_DIR := ./migrations

include .env
export DB_URL := $(DB_ADDR)

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help         Show this help message"
	@echo "  migrate      Create a new migration file (usage: make migrate name=<migration_name>)"
	@echo "  migrate-up   Run all pending migrations"
	@echo "  migrate-down Rollback the last migration"
	@echo "  migrate-force Force set migration version (usage: make migrate-force version=<version>)"
	@echo "  migrate-version Show current migration version"
	@echo "  dev          Start the development server with live reload"
	@echo "  infra-up      Start the infrastructure using Docker Compose"
	@echo "  seed         Seed the database with initial data"

dev:
	@air	

## Infrastructure targets
infra-up: ## Start the infrastructure using Docker Compose
	@docker-compose up 

infra-down: ## Stop the infrastructure using Docker Compose
	@docker-compose down	

## Documentation targets
gen-docs:
	@swag init -g ./cmd/api/main.go -d . -o swagger && swag fmt

## Migration targets	

migrate: ## Create a new migration file (usage: make migrate name=<migration_name>)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate name=<migration_name>"; exit 1; fi
	@migrate create -seq -ext sql -dir $(MIGRATION_DIR) $(name)

migrate-up: ## Run all pending migrations
	@migrate -path=$(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down: ## Rollback the last migration
	@migrate -path=$(MIGRATION_DIR) -database "$(DB_URL)" down 1

migrate-force: ## Force set migration version (usage: make migrate-force version=<version>)
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=<version>"; exit 1; fi
	@migrate -path=$(MIGRATION_DIR) -database "$(DB_URL)" force $(version)


migrate-version: ## Show current migration version
	@migrate -path=$(MIGRATION_DIR) -database "$(DB_URL)" version	

seed: ## Seed the database with initial data
	@go run ./cmd/seed/seed.go	
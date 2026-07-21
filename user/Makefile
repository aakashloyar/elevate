# Variables
# Load variables from .env
include .env
export

MIGRATIONS_DIR := migrations
MIGRATE_IMAGE := migrate/migrate:v4.19.1

DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: migrate-create migrate-up migrate-down

migrate-create:
	@read -p "Enter migration name: " name; \
	docker run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		$(MIGRATE_IMAGE) \
		create -ext sql -dir /migrations -seq $$name

migrate-up:
	docker run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		--network host \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(DB_URL)" \
		up

migrate-down:
	docker run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		--network host \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(DB_URL)" \
		down

hello:
	echo "Hello!"

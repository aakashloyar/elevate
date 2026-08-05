# Variables
# Load variables from .env
include .env
export

MIGRATIONS_DIR := migrations
MIGRATE_IMAGE := migrate/migrate:v4.19.1

DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

IMAGE_NAME := user-service
IMAGE_TAG := latest
REGISTRY := ghcr.io/aakashloyar

IMAGE := $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)


.PHONY: migrate-create migrate-up migrate-down

migrate-create:
	@read -p "Enter migration name: " name; \
	podman run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		$(MIGRATE_IMAGE) \
		create -ext sql -dir /migrations -seq $$name

migrate-up:
	podman run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		--network host \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(DB_URL)" \
		up

migrate-down:
	podman run --rm \
		-v $(shell pwd)/$(MIGRATIONS_DIR):/migrations \
		--network host \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(DB_URL)" \
		down

hello:
	echo "Hello!"

build-image:
	podman build -t $(IMAGE) .

push-image:
	podman push $(IMAGE)

pull-image:
	podman pull $(IMAGE)

run-image:
	podman run --rm -p 8080:8080 $(IMAGE)

clean-image:
	podman rmi $(IMAGE)

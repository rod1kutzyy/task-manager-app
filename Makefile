include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d notesapp-postgres

env-down:
	@docker compose down notesapp-postgres

env-cleanup:
	@read -p "Clear all volume files of the environment? The risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down notesapp-postgres && \
		rm -rf out/pgdata && \
		echo "The environment files are cleared"; \
	else \
		echo "Environment cleanup has been canceled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "The required seq parameter is missing. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm notesapp-postgres-migrate \
		create \
		-ext sql \
		-dir //migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "The required action parameter is missing. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm notesapp-postgres-migrate \
		-path //migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@notesapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
		go mod tidy && \
		go run cmd/notesapp/main.go

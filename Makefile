include .env
export

export PROJECT_ROOT=.


env-up:
	@docker compose up -d notesapp-postgres notesapp-redis

env-down:
	@docker compose down notesapp-postgres notesapp-redis

env-cleanup:
	@read -p "Clear all volume files of the environment? The risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down notesapp-postgres notesapp-redis port-forwarder web-server && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		rm -rf ${PROJECT_ROOT}/out/redis_data && \
		rm -rf ${PROJECT_ROOT}/out/caddy_data && \
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

logs-cleanup:
	@read -p "Clear all log files? The risk of logs loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "The log files are cleared"; \
	else \
		echo "Log files cleanup has been canceled"; \
	fi

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
		export POSTGRES_HOST=localhost && \
		export REDIS_HOST=localhost && \
		go mod tidy && \
		go run ${PROJECT_ROOT}/cmd/notesapp/main.go

app-deploy:
	@docker compose up -d --build notesapp

app-undeploy:
	@docker compose down notesapp

web-deploy:
	@docker compose up -d web-server

web-undeploy:
	@docker compose down web-server

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/notesapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps:
	@docker compose ps

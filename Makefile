ENV_FILE ?= .env.example

.PHONY: run up down logs

run:
	docker compose --env-file $(ENV_FILE) up --build --force-recreate

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f
ENV_FILE := .env.example
COMPOSE := docker compose

TOOLS_BIN := $(CURDIR)/bin
BUF := $(TOOLS_BIN)/buf

# Docker Compose configuration
.PHONY: run rebuild up down logs

run:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build

rebuild:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build --force-recreate

up:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build

down:
	$(COMPOSE) --env-file "$(ENV_FILE)" down

logs:
	$(COMPOSE) --env-file "$(ENV_FILE)" logs -f

# Protocol Buffers Compiler
.PHONY: tools proto proto-deps proto-lint

tools:
	GOBIN="$(TOOLS_BIN)" go install tool

proto: tools
	PATH="$(TOOLS_BIN):$$PATH" "$(BUF)" generate
	go mod tidy

proto-deps: tools
	"$(BUF)" dep update

proto-lint: tools
	"$(BUF)" lint

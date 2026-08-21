ENV_FILE := .env.example
COMPOSE := docker compose

TOOLS_BIN := $(CURDIR)/bin
BUF := $(TOOLS_BIN)/buf

# Управление Docker Compose
.PHONY: rebuild up down logs

rebuild:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build --force-recreate

up:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build

down:
	$(COMPOSE) --env-file "$(ENV_FILE)" down

logs:
	$(COMPOSE) --env-file "$(ENV_FILE)" logs -f

# Форматирование и тестирование Go-кода
.PHONY: fmt test check

fmt:
	go fmt ./pkg/...
	go -C services/catalog fmt ./...
	go -C services/gateway fmt ./...

test:
	go test ./pkg/...
	go -C services/catalog test ./...
	go -C services/gateway test ./...

check: fmt test

# Работа с Protocol Buffers
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

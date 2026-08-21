ENV_FILE := .env.example
COMPOSE := docker compose

TOOLS_BIN := $(CURDIR)/bin
TOOLS_MODULE := $(CURDIR)/tools
BUF := $(TOOLS_BIN)/buf
GO_PACKAGES := ./gen/... ./internal/... ./services/...

# Управление Docker Compose
.PHONY: run rebuild up down logs

run: rebuild

rebuild:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build --force-recreate

up:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build

down:
	$(COMPOSE) --env-file "$(ENV_FILE)" down

logs:
	$(COMPOSE) --env-file "$(ENV_FILE)" logs -f

# Проверка проекта
.PHONY: fmt test vet frontend-check check

fmt:
	GOWORK=off go fmt $(GO_PACKAGES)

test:
	GOWORK=off go test $(GO_PACKAGES)

vet:
	GOWORK=off go vet $(GO_PACKAGES)

frontend-check:
	npm --prefix frontend run lint
	npm --prefix frontend run build

check: fmt proto-lint test vet frontend-check

# Работа с Protocol Buffers
.PHONY: tools proto proto-deps proto-lint proto-breaking

tools:
	GOWORK=off GOBIN="$(TOOLS_BIN)" go -C "$(TOOLS_MODULE)" install tool

proto: tools
	PATH="$(TOOLS_BIN):$$PATH" "$(BUF)" generate
	go mod tidy

proto-deps: tools
	"$(BUF)" dep update

proto-lint: tools
	"$(BUF)" lint

# Проверяет обратную совместимость proto относительно текущей ветки master.
proto-breaking: tools
	"$(BUF)" breaking --against '.git#branch=master'

ENV_FILE := .env.example
COMPOSE := docker compose

TOOLS_BIN := $(CURDIR)/bin

PROTO_DIR := proto
PROTO_OUT := gen/go
PROTO_FILES := $(shell find $(PROTO_DIR) -type f -name '*.proto' 2>/dev/null | sort)

# Docker Compose configuration
.PHONY: run up down logs

run:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build --force-recreate

up:
	$(COMPOSE) --env-file "$(ENV_FILE)" up --build

down:
	$(COMPOSE) --env-file "$(ENV_FILE)" down

logs:
	$(COMPOSE) --env-file "$(ENV_FILE)" logs -f

# Protocol Buffers Compiler
.PHONY: tools proto require-proto-tools

tools:
	GOBIN="$(TOOLS_BIN)" go install tool

require-proto-tools:
	@if ! command -v protoc >/dev/null 2>&1; then \
		echo "Required tool is not installed: protoc"; \
		exit 1; \
	fi

proto: require-proto-tools tools
	@if [ -z "$(PROTO_FILES)" ]; then \
		echo "No .proto files found in $(PROTO_DIR)"; \
		exit 1; \
	fi
	@mkdir -p "$(PROTO_OUT)"
	PATH="$(TOOLS_BIN):$$PATH" protoc \
		--proto_path="$(PROTO_DIR)" \
		--go_out="$(PROTO_OUT)" \
		--go_opt=paths=source_relative \
		--go-grpc_out="$(PROTO_OUT)" \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	go mod tidy

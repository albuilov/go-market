# Инструмент для проверки стандартного gRPC Health API.
FROM ghcr.io/grpc-ecosystem/grpc-health-probe:v0.4.56 AS health-probe

FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY . .

RUN --mount=type=cache,id=go-market-modules,target=/go/pkg/mod \
    --mount=type=cache,id=go-market-build,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/catalog ./services/catalog/cmd

FROM alpine:3.22 AS runtime

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/catalog ./catalog
COPY --from=health-probe \
    /ko-app/grpc-health-probe \
    /usr/local/bin/grpc_health_probe

USER app

ENTRYPOINT ["./catalog"]

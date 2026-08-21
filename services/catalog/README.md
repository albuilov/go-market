# Сервис Catalog

Catalog владеет данными каталога и предоставляет внутренний gRPC API.
HTTP-маршруты для клиента описаны аннотациями в
`proto/catalog/v1/catalog.proto` и публикуются через Gateway.

## Конфигурация

- `CATALOG_GRPC_ADDRESS` — адрес gRPC-сервера, по умолчанию `:50051`.

## Локальная проверка

Из этого каталога сервис можно проверить без `go.work`:

```bash
GOWORK=off go test ./...
GOWORK=off go run ./cmd
```

При запуске вне Docker значение `CATALOG_GRPC_ADDRESS` можно не задавать.
Сервис публикует стандартный gRPC Health API для
`catalog.v1.CatalogService`.

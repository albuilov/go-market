# Архитектура

## Структура Go-кода

Runtime-код проекта находится в одном Go-модуле `go-market`. Это позволяет
сервисам переиспользовать инфраструктуру без зависимости от `go.work` и без
версирования внутренних модулей. Отдельный модуль `tools` содержит только
закреплённые версии генераторов.

```text
cmd сервиса -> internal/app -> transport и clients
                         \-> internal/platform
```

`internal/platform` содержит только компоненты, которые нужны нескольким
сервисам:

- `envconfig` — чтение переменных окружения;
- `logging` — единый структурированный JSON-логгер;
- `grpcserver` — interceptors, Protovalidate, Health и graceful shutdown;
- `grpcclient` — создание соединения и timeout interceptor;
- `health` — объединение readiness-проверок;
- `identity` — передача проверенной identity во внутренней gRPC metadata;
- `requestid` — общий ключ request ID.

Сервисная логика остаётся рядом с владельцем. Например, JWT verifier и policy
живут в Gateway, потому что другие сервисы не проверяют внешний токен.

## Поток запроса

1. Gateway принимает HTTP-запрос и назначает request ID.
2. grpc-gateway выбирает RPC по `google.api.http` аннотации.
3. Gateway проверяет JWT только для приватного RPC и передаёт проверенную
   identity во внутренней metadata.
4. gRPC-сервер выполняет request ID, logging, recovery и Protovalidate
   interceptors до бизнес-обработчика.
5. Ответ protobuf преобразуется в публичный JSON.

Порядок interceptors является частью поведения приложения. Внешние
interceptors видят результат внутренних и поэтому logging расположен до
recovery и validation.

## Зависимости сервисов

Gateway агрегирует readiness обязательных downstream-сервисов. Catalog
публикует стандартный gRPC Health API. При добавлении необязательной
зависимости её отказ не следует автоматически включать в readiness: сначала
нужно определить допустимую деградацию.

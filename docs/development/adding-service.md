# Добавление нового микросервиса

Новый сервис должен собираться и тестироваться из корневого Go-модуля без
`go.work`. Минимальный порядок действий:

1. Создать `services/<name>/cmd/main.go`, `internal/config` и `internal/app`.
2. Описать типизированный конфиг. Адрес сервера и timeout завершения работы
   должны быть явными полями.
3. В `main` загрузить конфиг, создать `logging.New`, контекст от системных
   сигналов и вызвать `app.Run`.
4. Описать protobuf-контракт в `proto/<name>/v1`, не переиспользуя номера
   удалённых полей.
5. Выполнить `make proto`, зарегистрировать сгенерированный gRPC service на
   `grpcserver.Server.GRPC()`.
6. После успешной инициализации вызвать `SetServing`, затем `Run`. Общий server
   уже включает Protovalidate, Health и bounded graceful shutdown.
7. Для downstream-зависимостей использовать `grpcclient.New`, явный timeout и
   подходящие transport credentials.
8. Добавить сервис в Docker Compose, healthcheck, `.env.example` и readiness
   его потребителей.
9. Добавить тесты конфига, transport и бизнес-логики, затем выполнить
   `make check`.

Пример каркаса инициализации:

```text
Load config
  -> create logger
  -> create gRPC server
  -> register service handlers
  -> mark services as serving
  -> run until context cancellation
```

Не следует переносить код в `internal/platform` заранее. Сначала компонент
должен появиться хотя бы в двух сервисах и иметь одинаковую семантику.

TODO: когда появится первый streaming RPC, добавить стандартную цепочку stream
interceptors рядом с unary-цепочкой.

# go-market

E-commerce backend на Go для практики и экспериментов с микросервисной архитектурой.

## Сервисы

- `frontend` — React SPA, собирается Vite и раздаётся через Caddy;
- `gateway` — внешняя HTTP-точка входа;
- `catalog` — управление каталогом товаров, предоставляет gRPC API.

## Требования

Для запуска приложения:

- Docker;
- Docker Compose;
- Make.

Для локальной разработки дополнительно:

- Go 1.26.4;
- Node.js 24;
- Protocol Buffers Compiler (`protoc`).

Проверить основные инструменты:

```bash
go version
node --version
npm --version
docker --version
docker compose version
protoc --version
make --version
```

### Установка компилятора Protocol Buffers на macOS

```bash
brew install protobuf
```

## Настройка окружения

Для локального запуска достаточно использовать `.env.example` файл конфигурации, в нем указаны значения по умолчанию.

## Локальный запуск

Собрать образы, пересоздать контейнеры и запустить сервисы с `.env.example`:

```bash
make run
```

После запуска доступны:

- frontend: [http://localhost:8080](http://localhost:8080);
- gateway API: [http://localhost:3000/api/v1/products](http://localhost:3000/api/v1/products);
- Swagger UI: [http://localhost:3001/docs/](http://localhost:3001/docs/);
- OpenAPI JSON: [http://localhost:3001/openapi.json](http://localhost:3001/openapi.json).

Swagger работает на отдельном debug-порту, который Docker Compose публикует только на `127.0.0.1`.

Frontend обращается к backend по относительному пути `/api`. В Docker запросы
проксирует Caddy на сервис `gateway`, поэтому отдельная настройка CORS не нужна.

Остановить локальное приложение:

```bash
make down
```

Посмотреть логи:

```bash
make logs
```

### Разработка frontend без Docker

Запустить backend через Docker Compose, затем в отдельном терминале:

```bash
cd frontend
npm install
npm run dev
```

Vite откроет приложение на `http://localhost:5173` и проксирует `/api` на
`http://localhost:3000`.

Проверить frontend перед сборкой контейнера:

```bash
cd frontend
npm run lint
npm run build
```

## Форматирование и тестирование Go-кода

Отформатировать Go-код во всех модулях:

```bash
make fmt
```

Запустить все Go-тесты:

```bash
make test
```

Последовательно выполнить форматирование и тесты:

```bash
make check
```

## Запуск с другим конфигурационным файлом

Для команд `up`, `down` и `logs` необходимо явно указать env-файл:

```bash
make up ENV_FILE=.env.local
make logs ENV_FILE=.env.local
make down ENV_FILE=.env.local
```

## Генерация protobuf

Исходные protobuf-контракты находятся в каталоге:

```text
proto/
```

Сгенерированный Go-код сохраняется в:

```text
gen/go/
```

Сгенерированная спецификация OpenAPI сохраняется в:

```text
gen/openapi/spec/go-market.swagger.json
```

Сгенерировать код для всех protobuf-контрактов:

```bash
make proto
```

Команда автоматически:

1. устанавливает закреплённые Go-инструменты в локальный каталог `bin`;
2. находит все `.proto` внутри каталога `proto`;
3. генерирует Go-структуры, gRPC- и gRPC-Gateway-код;
4. генерирует единую спецификацию OpenAPI 2.0;
5. обновляет зависимости корневого Go-модуля.

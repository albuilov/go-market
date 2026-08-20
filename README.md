# go-market

E-commerce backend на Go для практики и экспериментов с микросервисной архитектурой.

## Сервисы

- `gateway` — внешняя HTTP-точка входа;
- `catalog` — управление каталогом товаров, предоставляет gRPC API.

## Требования

Для запуска приложения:

- Docker;
- Docker Compose;
- Make.

Для локальной разработки дополнительно:

- Go 1.26.4;
- Protocol Buffers Compiler (`protoc`).

Проверить основные инструменты:

```bash
go version
docker --version
docker compose version
protoc --version
make --version
```

### Установка Protocol Buffers Compiler на macOS

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

Остановить локальное приложение:

```bash
make down
```

Посмотреть логи:

```bash
make logs
```

## Тестовый API

Получить моковый список товаров через Gateway:

```bash
curl http://localhost:8080/products
```

Gateway принимает HTTP-запрос, вызывает `CatalogService/ListProducts` по gRPC
и возвращает JSON:

```json
{
  "products": [
    {
      "id": "product-1",
      "name": "Mechanical Keyboard",
      "price_minor_units": 129900,
      "currency_code": "RUB"
    }
  ]
}
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

Сгенерировать код для всех protobuf-контрактов:

```bash
make proto
```

Команда автоматически:

1. устанавливает закреплённые Go-инструменты в локальный каталог `bin`;
2. находит все `.proto` внутри каталога `proto`;
3. генерирует Go-структуры и gRPC-код;
4. обновляет зависимости корневого Go-модуля.

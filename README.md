# URL Shortener

Сервис сокращения ссылок на Go.

## Запуск

### Docker Compose

```bash
docker-compose up --build
```

### Локально

**PostgreSQL:**

```bash
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=urlshortener -p 5432:5432 postgres:15
```

**Сервис:**

```bash
go run cmd/main.go -d "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"
```

## API

| Метод | Путь       | Описание                |
|-------|------------|-------------------------|
| POST  | `/`        | Создать короткую ссылку |
| GET   | `/{id}`    | Редирект                |
| GET   | `/ping`    | Health check            |

## Конфигурация

| Флаг | Env            | По умолчанию           |
|------|----------------|------------------------|
| `-a` | `SERVER_ADDRESS` | `localhost:8080`     |
| `-b` | `BASE_URL`       | `http://localhost:8080` |
| `-l` | `LOG_LEVEL`      | `info`                 |
| `-d` | `DATABASE_DSN`   | —                      |

## Тесты

```bash
go test ./tests/...
```

## Примеры

```bash
# POST запрос (создать короткую ссылку)
curl -X POST -d "https://example.com" http://localhost:8080
# Ответ: http://localhost:8080/a1b2c3d4

# GET запрос (перейти по ссылке)
curl -v http://localhost:8080/a1b2c3d4
# Редирект 307 на https://example.com

# Health check
curl http://localhost:8080/ping
# OK
```

## Команды Makefile

```bash
make run    # запустить сервер
make build  # собрать бинарник
make test   # запустить тесты
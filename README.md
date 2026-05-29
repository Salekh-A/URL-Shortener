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
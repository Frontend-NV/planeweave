# Planeweave — бекенд

REST API и WebSocket для курсового проекта Planeweave.

**Стек:** Gin, sqlc, PostgreSQL, gorilla/websocket.

## Архитектура (слои)

```
HTTP-запрос
  → handler/     — Gin-маршруты, JSON, коды ответов
  → service/     — бизнес-логика, проверки доступа, события WS
  → repository/  — доступ к данным (обёртка над sqlc)
  → db/queries/  — SQL-запросы (генерируются sqlc)
```

| Слой | Папка | Ответственность |
|------|-------|-----------------|
| Handler | `internal/handler/` | HTTP, привязка JSON, маршрутизация |
| Service | `internal/service/` | правила предметной области, bcrypt, DAG, broadcast |
| Repository | `internal/repository/` | транзакции, вызов sqlc |
| SQL | `db/queries/` + `internal/repository/sqlc/` | запросы к PostgreSQL |
| Domain | `internal/domain/` | модели и доменные ошибки |

## Быстрый старт

```bash
# из корня репозитория — PostgreSQL + API одной командой
make dev

# или по шагам:
make up          # PostgreSQL
make run         # API на :8080
```

Проверка: `curl http://localhost:8080/health`

**Swagger UI:** http://localhost:8080/swagger/index.html

Полный список команд: `make help`

## Swagger (swaggo)

Документация REST **генерируется из аннотаций** в `internal/handler/` и поднимается вместе с API:

| URL | Назначение |
|-----|------------|
| `/swagger/index.html` | интерактивный UI |
| `/swagger/doc.json` | OpenAPI JSON |

Обновить после изменения handler'ов:

```bash
make swag
```

Копия для фронта: `api/openapi.json`. Подробнее: `api/README.md`.

## sqlc

После изменения SQL в `db/queries/` или схемы в `db/migrations/`:

```bash
cd backend
sqlc generate
```

Конфигурация: `sqlc.yaml`.

## Документация

**Студентам (фронтенд):**

- **[Руководство по API](../docs/course/backend-dev-guide.md)** — как подключиться, JWT, curl, WS, ошибки
- Swagger UI: http://localhost:8080/swagger/index.html (при `make dev`)
- OpenAPI JSON: `api/openapi.json` (генерируется `make swag`)
- WebSocket: `docs/websocket.md`
- Product brief: `../docs/course/product-brief.md`

**Ментору:** архитектура и sqlc — ниже в этом файле.

## Переменные окружения

| Переменная   | По умолчанию |
|--------------|--------------|
| DATABASE_URL | postgres://planeweave:planeweave@localhost:5432/planeweave?sslmode=disable |
| JWT_SECRET   | dev-secret-change-me |
| HTTP_ADDR    | :8080 |
| CORS_ORIGINS | http://localhost:5173 |

## Тесты

```bash
cd backend
go test ./...
```

## График выкладки на staging

| Неделя курса | Модуль |
|--------------|--------|
| до 1 | brief, OpenAPI, WS doc |
| 3 | auth |
| 4 | projects, invites |
| 5 | чтение tasks, dependencies |
| 6 | CRUD + проверка DAG |
| 7 | WebSocket |

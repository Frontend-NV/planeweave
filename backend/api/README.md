# OpenAPI-спецификация

Спецификация **генерируется из кода** через [swaggo/swag](https://github.com/swaggo/swag) (аннотации в `internal/handler/`).

## Где смотреть

| Источник | URL / путь |
|----------|------------|
| **Swagger UI** (при запущенном API) | http://localhost:8080/swagger/index.html |
| JSON live | http://localhost:8080/swagger/doc.json |
| Файл (после `make swag`) | `api/openapi.json` |
| Исходники генерации | `docs/swag/swagger.yaml` |

## Обновить спецификацию

```bash
make swag
```

Запускается автоматически перед `make dev`, `make run` и `make build`.

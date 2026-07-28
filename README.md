# planeweave

Командная доска задач с графом зависимостей — курсовой проект.

## Документы

- [Спецификация / роадмап](docs/superpowers/specs/2026-07-28-planeweave-roadmap-design.md)
- [План реализации](docs/superpowers/plans/2026-07-28-planeweave-implementation-plan.md)
- [Product brief для студентов](docs/course/product-brief.md)
- **[Руководство: бекенд для разработки фронта](docs/course/backend-dev-guide.md)** ← студентам
- [Бекенд — README (ментор)](backend/README.md)

## Структура

| Папка | Кто пишет | Содержимое |
|-------|-----------|------------|
| `backend/` | ментор | Go API, OpenAPI, WebSocket |
| `docs/course/` | ментор + студенты | brief, шаблон ТЗ, макеты команд |
| `frontend/` | **студенты** | React-клиент (создаётся на неделе 3) |

Локальный запуск: `make dev` (см. `Makefile` в корне).

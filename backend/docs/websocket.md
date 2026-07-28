# Канал синхронизации Planeweave

Документ описывает WebSocket-протокол для синхронизации доски проекта в реальном времени.

## Подключение

```
ws://{host}/ws?token={jwt}&projectId={uuid}
wss://{host}/ws?token={jwt}&projectId={uuid}   # production
```

| Параметр    | Обязателен | Описание                                      |
|-------------|------------|-----------------------------------------------|
| `token`     | да         | JWT из `/auth/login` или `/auth/register`     |
| `projectId` | да         | UUID проекта, к которому нужен доступ         |

**Проверки при подключении:**
- токен валиден;
- пользователь — участник проекта.

При успехе сервер отправляет:

```json
{
  "type": "connected",
  "payload": {
    "projectId": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## Формат сообщений

Все сообщения — JSON с полями:

```json
{
  "type": "event.name",
  "payload": { }
}
```

Клиент **не отправляет** мутации через WebSocket. Изменения — только через REST; сервер рассылает события участникам комнаты.

## События сервера → клиент

### `task.created`

Новая задача в проекте.

```json
{
  "type": "task.created",
  "payload": {
    "id": "uuid",
    "projectId": "uuid",
    "title": "Настроить CI",
    "status": "todo",
    "assigneeId": null,
    "x": 120,
    "y": 80,
    "updatedAt": "2026-07-28T12:00:00Z"
  }
}
```

### `task.updated`

Изменены поля задачи (статус, название, исполнитель).

```json
{
  "type": "task.updated",
  "payload": {
    "id": "uuid",
    "projectId": "uuid",
    "title": "Настроить CI",
    "status": "in_progress",
    "assigneeId": "uuid",
    "x": 120,
    "y": 80,
    "updatedAt": "2026-07-28T12:05:00Z"
  }
}
```

### `task.deleted`

```json
{
  "type": "task.deleted",
  "payload": {
    "id": "uuid",
    "projectId": "uuid"
  }
}
```

### `task.moved`

Только позиция на канвасе (можно объединять с `task.updated` — клиент обрабатывает оба одинаково).

```json
{
  "type": "task.moved",
  "payload": {
    "id": "uuid",
    "projectId": "uuid",
    "x": 200,
    "y": 150,
    "updatedAt": "2026-07-28T12:06:00Z"
  }
}
```

### `dependency.created`

```json
{
  "type": "dependency.created",
  "payload": {
    "id": "uuid",
    "projectId": "uuid",
    "fromTaskId": "uuid",
    "toTaskId": "uuid"
  }
}
```

### `dependency.deleted`

```json
{
  "type": "dependency.deleted",
  "payload": {
    "id": "uuid",
    "projectId": "uuid"
  }
}
```

### `presence.updated`

Количество участников онлайн в комнате проекта.

```json
{
  "type": "presence.updated",
  "payload": {
    "projectId": "uuid",
    "onlineCount": 2
  }
}
```

## Рекомендации для клиента

1. После REST-мутации инициатор может получить то же событие — сравнивайте `updatedAt`, чтобы не дублировать оптимистичное обновление.
2. При разрыве соединения показывайте плашку «Переподключение…» и переподключайтесь с задержкой 1с → 2с → 4с → … (максимум 30с).
3. После восстановления канала выполните полную загрузку задач и зависимостей через REST.
4. REST остаётся рабочим, если WebSocket недоступен.

## Коды ошибок REST (связанные с графом)

| HTTP | `error`           | Когда                          |
|------|-------------------|--------------------------------|
| 409  | `cycle_detected`  | новая зависимость создаёт цикл |

## Готовность по неделям курса

| Неделя | Что должно работать на стенде      |
|--------|------------------------------------|
| 7      | все события выше + `presence.updated` |

До недели 7 клиент может работать только через REST.

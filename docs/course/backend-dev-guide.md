# Руководство студента: работа с бекендом Planeweave

Бекенд пишет и поднимает ментор. **Студенты реализуют только фронтенд** и ходят в API по контракту. Этот документ — как подключиться, отладить клиент и что ожидать от сервера.

---

## 1. Что вам нужно знать

| Документ | Зачем |
|----------|--------|
| **Swagger UI** (при запущенном API) | http://localhost:8080/swagger/index.html — интерактивная документация |
| [`backend/api/openapi.json`](../../backend/api/openapi.json) | OpenAPI JSON (генерируется из кода, `make swag`) |
| [`backend/docs/websocket.md`](../../backend/docs/websocket.md) | синхронизация в реальном времени |
| [`product-brief.md`](product-brief.md) | суть продукта (для ТЗ на неделе 1) |
| [`tz-template.md`](tz-template.md) | шаблон ТЗ; образец структуры — [`ВКР ТЗ.pdf`](../examples/ВКР%20ТЗ.pdf) |

Код бекенда (`backend/internal/…`) **можно не читать** — ориентируйтесь на Swagger и этот гайд.

---

## 2. Локальный API (если ментор не дал staging)

Ментор может дать URL тестового стенда. Если работаете локально — из **корня репозитория**:

```bash
make dev
```

Поднимется PostgreSQL и API на `http://localhost:8080`.

**Swagger UI:** http://localhost:8080/swagger/index.html

Проверка:

```bash
curl http://localhost:8080/health
# ok
```

Остановить только базу: `make down`. Полный список: `make help`.

**Важно:** бекенд локально поднимает ментор или тот, у кого установлены Docker и Go. Достаточно, чтобы **хотя бы один человек в команде** мог запустить API для остальных.

---

## 3. Переменные окружения фронтенда

В папке `frontend/` (создаёте на неделе 3):

```env
# .env.local
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

Если ментор дал staging:

```env
VITE_API_URL=https://api.staging.example.com
VITE_WS_URL=wss://api.staging.example.com/ws
```

CORS на бекенде по умолчанию разрешает `http://localhost:5173` (Vite). Если порт другой — попросите ментора добавить origin в `CORS_ORIGINS`.

---

## 4. Авторизация (JWT)

После регистрации или входа сервер возвращает:

```json
{
  "token": "eyJhbG...",
  "user": {
    "id": "uuid",
    "email": "you@example.com",
    "displayName": "Имя"
  }
}
```

**Все запросы к API** (кроме `/auth/register`, `/auth/login`, `/health`) — с заголовком:

```
Authorization: Bearer <token>
```

Рекомендации для фронта:

1. Сохранить `token` (часто `localStorage` — как договоритесь в команде).
2. При загрузке приложения — `GET /auth/me` для проверки сессии.
3. При ответе **401** — очистить токен и перенаправить на страницу входа.

### Пример: регистрация и вход

```bash
# регистрация
curl -s -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "student@test.com",
    "password": "password1",
    "displayName": "Студент"
  }'

# вход
curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"student@test.com","password":"password1"}'
```

Сохраните `token` из ответа:

```bash
export TOKEN="eyJ..."
```

---

## 5. Типы и клиент из OpenAPI

Спецификация генерируется из Go-кода (swaggo). Актуальный JSON:

- в репозитории: `backend/api/openapi.json` (после `make swag`);
- с запущенного API: http://localhost:8080/swagger/doc.json .

```bash
cd frontend
npm install -D openapi-typescript
npx openapi-typescript ../backend/api/openapi.json -o src/shared/api/schema.d.ts
```

Минимальный HTTP-клиент:

```typescript
const baseUrl = import.meta.env.VITE_API_URL;

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token');
  const res = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw { status: res.status, ...body };
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}
```

---

## 6. Типичный сценарий разработки (REST)

Порядок, близкий к неделям курса:

### Неделя 3–4: auth и проекты

```bash
# список проектов
curl -s http://localhost:8080/projects -H "Authorization: Bearer $TOKEN"

# создать проект
curl -s -X POST http://localhost:8080/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Мой проект"}'
# запомните id проекта → PROJECT_ID
```

### Неделя 4: приглашение второго участника

```bash
# создать invite (от имени участника 1)
curl -s -X POST "http://localhost:8080/projects/$PROJECT_ID/invites" \
  -H "Authorization: Bearer $TOKEN"
# в ответе token — часть ссылки /invite/:token

# участник 2 (после своего login):
curl -s -X POST "http://localhost:8080/invites/$INVITE_TOKEN/accept" \
  -H "Authorization: Bearer $TOKEN2"
```

### Неделя 5–6: задачи и зависимости

```bash
# задачи проекта
curl -s "http://localhost:8080/projects/$PROJECT_ID/tasks" \
  -H "Authorization: Bearer $TOKEN"

# создать задачу
curl -s -X POST "http://localhost:8080/projects/$PROJECT_ID/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Задача A","x":100,"y":100}'

# зависимость: A блокирует B (from → to)
curl -s -X POST "http://localhost:8080/projects/$PROJECT_ID/dependencies" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"fromTaskId":"UUID_A","toTaskId":"UUID_B"}'

# изменить задачу
curl -s -X PATCH "http://localhost:8080/tasks/$TASK_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"status":"in_progress"}'
```

**Статусы задачи:** `todo`, `in_progress`, `done`.

---

## 7. Ошибки API (что показывать пользователю)

| HTTP | `error` в JSON | Действие на фронте |
|------|----------------|---------------------|
| 401 | `unauthorized`, `invalid_credentials` | на страницу входа |
| 403 | `forbidden` | «Нет доступа к проекту» |
| 404 | `not_found`, `invite_not_found` | сообщение + возврат к списку |
| 409 | `cycle_detected` | «Циклическая зависимость запрещена», откат связи |
| 409 | `email_taken` | «Email уже занят» |
| 422 | `invalid_request`, `password_too_short` | подсветка полей формы |
| 500 | `internal_error` | общее сообщение об ошибке |

Формат ошибки:

```json
{ "error": "cycle_detected" }
```

Проверку цикла на клиенте делайте **до** запроса (см. план курса), но сервер всё равно вернёт 409 при цикле.

---

## 8. WebSocket (неделя 7)

Подробности — в [`backend/docs/websocket.md`](../../backend/docs/websocket.md).

**Подключение:**

```
ws://localhost:8080/ws?token=<JWT>&projectId=<UUID проекта>
```

**Правило:** изменения отправляете через **REST**; по каналу приходят только **события** для обновления UI у всех участников.

События: `task.created`, `task.updated`, `task.deleted`, `task.moved`, `dependency.created`, `dependency.deleted`, `presence.updated`.

Пример обработки:

```typescript
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  switch (msg.type) {
    case 'task.updated':
      queryClient.setQueryData(['tasks', projectId], (old) => /* обновить */);
      break;
    // ...
  }
};
```

При разрыве соединения — переподключение с задержкой; после восстановления — полная загрузка задач через REST.

---

## 9. Отладка без фронта

Полезно проверить API до React:

1. **curl** — примеры выше.
2. **Swagger UI** — http://localhost:8080/swagger/index.html (кнопка «Authorize», вставьте `Bearer <token>`).
3. **Postman / Insomnia** — импорт `backend/api/openapi.json` или URL `http://localhost:8080/swagger/doc.json`.
4. **Два браузера / два профиля** — два аккаунта, один проект через invite, проверка WS на неделе 7.

Проверка health:

```bash
curl http://localhost:8080/health
```

---

## 10. Частые проблемы

| Симптом | Решение |
|---------|---------|
| `Failed to fetch` / CORS | фронт на `localhost:5173`? иначе — ментор добавляет origin |
| 401 на все запросы | проверьте заголовок `Authorization: Bearer …` |
| пустой список проектов | создайте проект через POST `/projects` |
| invite не работает | ссылка `/invite/:token`, второй пользователь должен быть залогинен |
| WS не подключается | `token` и `projectId` в query; пользователь — участник проекта |
| API не отвечает | `make up` или спросите URL staging у ментора |

---

## 11. Что не нужно делать студентам

- Менять `backend/` без согласования с ментором.
- Дублировать бизнес-правила (DAG, права доступа) только на клиенте — сервер всё равно проверяет.
- Отправлять мутации через WebSocket — только REST.

---

## 12. Связь с неделями курса

| Неделя | REST | WS |
|--------|------|-----|
| 3 | register, login, me | — |
| 4 | projects, invites | — |
| 5 | GET tasks, GET dependencies | — |
| 6 | POST/PATCH/DELETE tasks и dependencies | — |
| 7 | как раньше | подключение к `/ws` |

Если на вашей неделе модуль на staging ещё не включён — ментор сообщит; можно временно мокать ответы или договориться о локальном `make dev`.

---

**Вопросы по API** — к ментору, с указанием endpoint, тела запроса и кода ответа.

# Planeweave — план реализации курсового проекта

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Реализовать Planeweave — командную доску задач с графом зависимостей и синхронизацией в реальном времени: бекенд на Go (ментор) и фронтенд на React (студенты) за 10 недель.

**Architecture:** REST API для CRUD и JWT-авторизации; WebSocket-комната на проект для рассылки событий; фронтенд хранит серверные данные в TanStack Query, локальный UI — в Zustand, канвас — ReactFlow. Проверка ацикличности графа — на бекенде (409) и дублируется на клиенте для мгновенной обратной связи.

**Tech Stack:** Go 1.22+, PostgreSQL, gorilla/websocket или nhooyr.io/websocket; TypeScript, Vite, React 18, React Router, TanStack Query, Zustand, @xyflow/react, Vitest.

**Спецификация:** `docs/superpowers/specs/2026-07-28-planeweave-roadmap-design.md`

---

## Карта файлов

### Бекенд (ментор)

```
backend/
├── cmd/server/main.go                 # точка входа, маршруты, WS
├── api/openapi.yaml                   # контракт для студентов
├── docs/websocket.md                  # события канала связи
├── internal/
│   ├── config/config.go
│   ├── auth/handler.go, service.go, jwt.go
│   ├── projects/handler.go, service.go
│   ├── tasks/handler.go, service.go
│   ├── dependencies/handler.go, service.go, cycle.go
│   ├── invites/handler.go, service.go
│   ├── ws/hub.go, client.go, events.go
│   └── store/postgres.go, migrations/
├── migrations/001_init.sql
├── go.mod
└── Dockerfile
```

### Фронтенд (студенты)

```
frontend/
├── src/
│   ├── app/App.tsx, routes.tsx, providers.tsx
│   ├── pages/auth/LoginPage.tsx, RegisterPage.tsx
│   ├── pages/projects/ProjectListPage.tsx, ProjectSettingsPage.tsx
│   ├── pages/board/BoardPage.tsx
│   ├── pages/invite/InvitePage.tsx
│   ├── features/auth/store.ts, hooks.ts, api.ts
│   ├── features/projects/api.ts, hooks.ts
│   ├── features/tasks/api.ts, hooks.ts, types.ts
│   ├── features/dependencies/api.ts, cycleCheck.ts
│   ├── features/realtime/wsClient.ts, eventHandlers.ts
│   ├── widgets/board/BoardCanvas.tsx, TaskNode.tsx, DependencyEdge.tsx
│   ├── widgets/board/TaskPanel.tsx
│   ├── shared/api/client.ts, generated/schema.d.ts
│   ├── shared/ui/Button.tsx, Input.tsx, Toast.tsx, Spinner.tsx
│   └── shared/lib/errors.ts, constants.ts
├── tests/features/dependencies/cycleCheck.test.ts
├── tests/pages/LoginPage.test.tsx
├── tests/pages/ProjectListPage.test.tsx
├── .env.example
├── vite.config.ts
└── package.json
```

### Документы курса (студенты, недели 1–2)

```
docs/course/
├── product-brief.md          # выдаёт ментор
├── tz-template.md            # шаблон ТЗ
└── design/                   # макеты команды
```

---

## Часть A. Подготовка ментора (до недели 1)

### Task A1: Product brief

**Files:**
- Create: `docs/course/product-brief.md`

- [ ] **Шаг 1:** Написать brief (1–2 страницы): проблема, пользователи, сценарии, ограничения DAG, real-time
- [ ] **Шаг 2:** Приложить скетч экранов (можно ASCII или ссылку на Figma ментора)
- [ ] **Шаг 3:** Раздать командам до старта недели 1

---

### Task A2: OpenAPI-контракт

**Files:**
- Create: `backend/api/openapi.yaml`

- [ ] **Шаг 1:** Описать схемы

```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id: { type: string, format: uuid }
        email: { type: string, format: email }
        displayName: { type: string }
    Task:
      type: object
      properties:
        id: { type: string, format: uuid }
        projectId: { type: string, format: uuid }
        title: { type: string }
        status: { type: string, enum: [todo, in_progress, done] }
        assigneeId: { type: string, format: uuid, nullable: true }
        x: { type: number }
        y: { type: number }
        updatedAt: { type: string, format: date-time }
    Dependency:
      type: object
      properties:
        id: { type: string, format: uuid }
        projectId: { type: string, format: uuid }
        fromTaskId: { type: string, format: uuid }
        toTaskId: { type: string, format: uuid }
```

- [ ] **Шаг 2:** Описать конечные точки (минимальный набор)

| Метод | Путь | Назначение |
|-------|------|------------|
| POST | `/auth/register` | регистрация |
| POST | `/auth/login` | вход → JWT |
| GET | `/auth/me` | текущий пользователь |
| GET | `/projects` | список проектов |
| POST | `/projects` | создание |
| PATCH | `/projects/{id}` | переименование |
| DELETE | `/projects/{id}` | удаление |
| GET | `/projects/{id}/members` | участники |
| POST | `/projects/{id}/invites` | создать приглашение |
| POST | `/invites/{token}/accept` | принять |
| GET | `/projects/{id}/tasks` | задачи |
| POST | `/projects/{id}/tasks` | создать задачу |
| PATCH | `/tasks/{id}` | изменить |
| DELETE | `/tasks/{id}` | удалить |
| GET | `/projects/{id}/dependencies` | связи |
| POST | `/projects/{id}/dependencies` | создать (409 при цикле) |
| DELETE | `/dependencies/{id}` | удалить |

- [ ] **Шаг 3:** Опубликовать YAML в репозитории курса

---

### Task A3: Документ WebSocket

**Files:**
- Create: `backend/docs/websocket.md`

- [ ] **Шаг 1:** Описать подключение

```
URL: wss://{host}/ws?token={jwt}&projectId={uuid}
Первое сообщение сервера: { "type": "connected", "projectId": "..." }
Клиент подписан на комнату projectId автоматически
```

- [ ] **Шаг 2:** Описать события

```json
{ "type": "task.created", "payload": { /* Task */ } }
{ "type": "task.updated", "payload": { /* Task */ } }
{ "type": "task.deleted", "payload": { "id": "uuid" } }
{ "type": "task.moved", "payload": { "id": "uuid", "x": 0, "y": 0 } }
{ "type": "dependency.created", "payload": { /* Dependency */ } }
{ "type": "dependency.deleted", "payload": { "id": "uuid" } }
{ "type": "presence.updated", "payload": { "onlineCount": 2 } }
```

---

### Task A4: Бекенд — неделя 3 (auth)

**Files:**
- Create: `backend/migrations/001_init.sql`, `backend/internal/auth/*`, `backend/cmd/server/main.go`

- [ ] **Шаг 1:** Миграция users

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  display_name TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
```

- [ ] **Шаг 2:** Реализовать register/login/me с bcrypt + JWT
- [ ] **Шаг 3:** Развернуть staging, выдать `VITE_API_URL` командам
- [ ] **Шаг 4:** Проверить через curl регистрацию и вход

---

### Task A5: Бекенд — неделя 4 (проекты и приглашения)

**Files:**
- Create: `backend/migrations/002_projects.sql`, `backend/internal/projects/*`, `backend/internal/invites/*`

- [ ] **Шаг 1:** Таблицы projects, project_members, invites
- [ ] **Шаг 2:** CRUD проектов + проверка членства
- [ ] **Шаг 3:** POST invite → token; POST accept → добавить в project_members

---

### Task A6: Бекенд — неделя 5–6 (задачи, зависимости, DAG)

**Files:**
- Create: `backend/internal/tasks/*`, `backend/internal/dependencies/cycle.go`

- [ ] **Шаг 1:** Таблицы tasks, dependencies
- [ ] **Шаг 2:** CRUD задач с полями x, y, status, assignee_id
- [ ] **Шаг 3:** Реализовать проверку цикла

```go
// internal/dependencies/cycle.go
func WouldCreateCycle(deps []Dependency, fromID, toID string) bool {
    adj := buildAdjacency(deps)
    adj[fromID] = append(adj[fromID], toID)
    return isReachable(adj, toID, fromID)
}
```

- [ ] **Шаг 4:** POST dependency возвращает 409 + `{ "error": "cycle_detected" }` при цикле

---

### Task A7: Бекенд — неделя 7 (WebSocket)

**Files:**
- Create: `backend/internal/ws/hub.go`, `client.go`, `events.go`

- [ ] **Шаг 1:** Hub: map[projectID]→clients, broadcast при мутациях
- [ ] **Шаг 2:** После каждого PATCH task — `task.updated` в комнату
- [ ] **Шаг 3:** presence.updated при connect/disconnect
- [ ] **Шаг 4:** Выдать `VITE_WS_URL` командам

---

## Часть B. Неделя 1 — Техническое задание (студенты)

### Task B1: Техническое задание

**Files:**
- Шаблон: `docs/course/tz-template.md` (ментор)
- Образец оформления: `docs/examples/ВКР ТЗ.pdf` (ориентир, не эталон)
- Create: `docs/course/teams/{team-name}/tz.md` (студенты)

**Формат:** разделы 1–7 шаблона (как в примере ВКР, но короче). Не писать: ТЭО, правовые основания, Docker-маркировку, нагрузочные испытания.

- [ ] **Шаг 1:** Скопировать `tz-template.md` в папку команды
- [ ] **Шаг 2:** §1–2 — введение, актуальность, цели и роли (из `product-brief.md`)
- [ ] **Шаг 3:** §3 — кратко: что делает UI и что берётся из API ментора
- [ ] **Шаг 4:** §4.1.1 — 8–10 пользовательских историй с критериями (auth, проекты, доска, зависимости, синхронизация)
- [ ] **Шаг 5:** §4.1.2–4.1.3 — экраны, переходы, диаграмма сущностей
- [ ] **Шаг 6:** §4.1.4–4.1.5 — таблица REST-точек (OpenAPI) и WebSocket-событий (`backend/docs/websocket.md`)
- [ ] **Шаг 7:** §4.2–6 — нефункциональные требования, критерии приёмки, стадии по неделям курса
- [ ] **Шаг 8:** Pull request → ревью ментора → правки → label `tz-approved`

**Зачёт:** ментор ставит label `tz-approved` на PR.

---

## Часть C. Неделя 2 — Макеты (студенты)

### Task C1: Каркасы экранов

**Files:**
- Create: `docs/course/teams/{team-name}/design/*.png` (или ссылка на Figma)

- [ ] **Шаг 1:** Макет входа и регистрации
- [ ] **Шаг 2:** Макет списка проектов (пустое состояние + список)
- [ ] **Шаг 3:** Макет доски: канвас, боковая панель задачи, шапка с участниками
- [ ] **Шаг 4:** Макет настроек проекта (участники, кнопка «Пригласить»)
- [ ] **Шаг 5:** Макет принятия приглашения
- [ ] **Шаг 6:** Схема переходов: регистрация → проект → invite → доска
- [ ] **Шаг 7:** PR → ревью ментора → `design-approved`

---

## Часть D. Неделя 3 — Каркас и авторизация (студенты)

### Task D1: Инициализация фронтенда

**Files:**
- Create: `frontend/package.json`, `frontend/vite.config.ts`, `frontend/tsconfig.json`

- [ ] **Шаг 1:** Создать проект

```bash
cd frontend
npm create vite@latest . -- --template react-ts
npm install react-router-dom @tanstack/react-query zustand
npm install -D eslint prettier vitest @testing-library/react @testing-library/jest-dom jsdom
```

- [ ] **Шаг 2:** Создать `.env.example`

```
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

- [ ] **Шаг 3:** Настроить алиас `@/` → `src/` в `vite.config.ts`
- [ ] **Шаг 4:** Создать структуру папок из карты файлов

---

### Task D2: HTTP-клиент

**Files:**
- Create: `frontend/src/shared/api/client.ts`
- Create: `frontend/src/shared/api/generated/schema.d.ts` (из OpenAPI ментора)

- [ ] **Шаг 1:** Сгенерировать типы

```bash
npx openapi-typescript ../backend/api/openapi.yaml -o src/shared/api/generated/schema.d.ts
```

- [ ] **Шаг 2:** Реализовать клиент

```typescript
// frontend/src/shared/api/client.ts
const baseUrl = import.meta.env.VITE_API_URL;

export async function api<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
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

### Task D3: Хранилище авторизации

**Files:**
- Create: `frontend/src/features/auth/store.ts`

- [ ] **Шаг 1:** Zustand store

```typescript
import { create } from 'zustand';

type User = { id: string; email: string; displayName: string };

type AuthState = {
  user: User | null;
  token: string | null;
  setSession: (user: User, token: string) => void;
  clearSession: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem('token'),
  setSession: (user, token) => {
    localStorage.setItem('token', token);
    set({ user, token });
  },
  clearSession: () => {
    localStorage.removeItem('token');
    set({ user: null, token: null });
  },
}));
```

---

### Task D4: Страница входа (TDD)

**Files:**
- Create: `frontend/tests/pages/LoginPage.test.tsx`
- Create: `frontend/src/pages/auth/LoginPage.tsx`
- Create: `frontend/src/features/auth/api.ts`

- [ ] **Шаг 1:** Написать падающий тест

```typescript
// frontend/tests/pages/LoginPage.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LoginPage } from '@/pages/auth/LoginPage';

it('показывает ошибку при пустом email', async () => {
  render(<LoginPage />);
  await userEvent.click(screen.getByRole('button', { name: 'Войти' }));
  expect(await screen.findByText('Введите email')).toBeInTheDocument();
});
```

- [ ] **Шаг 2:** Запустить — ожидаем FAIL

```bash
npm run test -- LoginPage.test.tsx
```

- [ ] **Шаг 3:** Реализовать форму с zod-валидацией (email, password обязательны)
- [ ] **Шаг 4:** `login()` → POST `/auth/login` → `setSession` → redirect `/projects`
- [ ] **Шаг 5:** Тест PASS
- [ ] **Шаг 6:** Аналогично — `RegisterPage.tsx` (без TDD или один тест на обязательные поля)

---

### Task D5: Защищённые маршруты

**Files:**
- Create: `frontend/src/app/routes.tsx`, `frontend/src/app/providers.tsx`

- [ ] **Шаг 1:** Маршруты: `/login`, `/register`, `/projects`, `/projects/:id`, `/projects/:id/settings`, `/invite/:token`
- [ ] **Шаг 2:** `ProtectedRoute` — без token redirect на `/login`
- [ ] **Шаг 3:** При загрузке приложения — GET `/auth/me` для восстановления сессии
- [ ] **Шаг 4:** Кнопка «Выйти» → `clearSession` → `/login`

**Зачёт недели 3:** регистрация, вход, редирект на `/projects`.

---

## Часть E. Неделя 4 — Проекты и приглашения

### Task E1: Список проектов

**Files:**
- Create: `frontend/src/features/projects/api.ts`, `hooks.ts`
- Create: `frontend/src/pages/projects/ProjectListPage.tsx`
- Test: `frontend/tests/pages/ProjectListPage.test.tsx`

- [ ] **Шаг 1:** Тест: пустой список → текст «Создайте первый проект»
- [ ] **Шаг 2:** `useProjects()` — TanStack Query GET `/projects`
- [ ] **Шаг 3:** Карточки проектов, клик → `/projects/:id`
- [ ] **Шаг 4:** Модальное окно «Создать проект» → POST `/projects`
- [ ] **Шаг 5:** Тест PASS

---

### Task E2: Настройки и приглашения

**Files:**
- Create: `frontend/src/pages/projects/ProjectSettingsPage.tsx`
- Create: `frontend/src/pages/invite/InvitePage.tsx`

- [ ] **Шаг 1:** GET `/projects/:id/members` — список участников
- [ ] **Шаг 2:** POST `/projects/:id/invites` — показать ссылку `{origin}/invite/{token}`
- [ ] **Шаг 3:** `InvitePage` — POST `/invites/:token/accept` → redirect на доску
- [ ] **Шаг 4:** PATCH/DELETE проекта на странице настроек

**Зачёт недели 4:** два аккаунта в одном проекте через ссылку.

---

## Часть F. Неделя 5 — Граф (чтение)

### Task F1: Преобразование данных → ReactFlow

**Files:**
- Create: `frontend/src/widgets/board/mapToFlow.ts`

- [ ] **Шаг 1:** Функция маппинга

```typescript
import type { Node, Edge } from '@xyflow/react';
import type { Task, Dependency } from '@/features/tasks/types';

export function mapToFlow(tasks: Task[], deps: Dependency[]): { nodes: Node[]; edges: Edge[] } {
  const nodes: Node[] = tasks.map((t) => ({
    id: t.id,
    type: 'task',
    position: { x: t.x, y: t.y },
    data: { title: t.title, status: t.status, assigneeId: t.assigneeId },
  }));
  const edges: Edge[] = deps.map((d) => ({
    id: d.id,
    source: d.fromTaskId,
    target: d.toTaskId,
    type: 'dependency',
  }));
  return { nodes, edges };
}
```

---

### Task F2: Узел и ребро

**Files:**
- Create: `frontend/src/widgets/board/TaskNode.tsx`, `DependencyEdge.tsx`, `BoardCanvas.tsx`

- [ ] **Шаг 1:** `npm install @xyflow/react`
- [ ] **Шаг 2:** `TaskNode` — div с title, цвет по status, имя исполнителя
- [ ] **Шаг 3:** `DependencyEdge` — стандартная стрелка ReactFlow
- [ ] **Шаг 3:** `BoardCanvas` — ReactFlow с `nodeTypes`, `edgeTypes`, Controls, Background
- [ ] **Шаг 4:** `BoardPage` — `useTasks(projectId)`, `useDependencies(projectId)` → `mapToFlow`
- [ ] **Шаг 5:** Клик по узлу → Zustand `selectedTaskId` → `TaskPanel` (read-only)

**Зачёт недели 5:** доска отображает задачи и связи из API.

---

## Часть G. Неделя 6 — Редактирование и DAG

### Task G1: Проверка цикла (TDD)

**Files:**
- Create: `frontend/src/features/dependencies/cycleCheck.ts`
- Test: `frontend/tests/features/dependencies/cycleCheck.test.ts`

- [ ] **Шаг 1:** Тесты

```typescript
import { wouldCreateCycle } from '@/features/dependencies/cycleCheck';

const deps = [
  { fromTaskId: 'a', toTaskId: 'b' },
  { fromTaskId: 'b', toTaskId: 'c' },
];

it('запрещает цикл a→c при связи c→a', () => {
  expect(wouldCreateCycle(deps, 'c', 'a')).toBe(true);
});

it('разрешает связь d→a без цикла', () => {
  expect(wouldCreateCycle(deps, 'd', 'a')).toBe(false);
});

it('запрещает самосвязь', () => {
  expect(wouldCreateCycle(deps, 'a', 'a')).toBe(true);
});
```

- [ ] **Шаг 2:** RUN — FAIL
- [ ] **Шаг 3:** Реализация DFS от `toId` к `fromId`
- [ ] **Шаг 4:** RUN — PASS

---

### Task G2: CRUD на доске

**Files:**
- Modify: `frontend/src/widgets/board/BoardCanvas.tsx`
- Create: `frontend/src/widgets/board/TaskPanel.tsx`

- [ ] **Шаг 1:** Кнопка «Добавить задачу» → POST с x,y центра viewport
- [ ] **Шаг 2:** `TaskPanel` — редактирование title, status (select), assignee (select из members)
- [ ] **Шаг 3:** `onNodeDragStop` → PATCH `/tasks/:id` { x, y }
- [ ] **Шаг 4:** `onConnect` → проверка `wouldCreateCycle` → POST dependency или toast «Циклическая зависимость»
- [ ] **Шаг 5:** Удаление задачи/ребра — DELETE + обновление Query cache
- [ ] **Шаг 6:** Optimistic updates в мутациях TanStack Query

**Зачёт недели 6:** 5+ задач, зависимости, цикл блокируется.

---

## Часть H. Неделя 7 — WebSocket

### Task H1: Клиент канала связи

**Files:**
- Create: `frontend/src/features/realtime/wsClient.ts`, `eventHandlers.ts`

- [ ] **Шаг 1:** Класс/хук WebSocket

```typescript
export function connectProjectWs(projectId: string, token: string, onEvent: (e: WsEvent) => void) {
  const url = `${import.meta.env.VITE_WS_URL}?token=${token}&projectId=${projectId}`;
  let ws: WebSocket | null = null;
  let delay = 1000;

  function connect() {
    ws = new WebSocket(url);
    ws.onmessage = (msg) => onEvent(JSON.parse(msg.data));
    ws.onclose = () => setTimeout(connect, Math.min(delay *= 2, 30000));
  }
  connect();
  return () => ws?.close();
}
```

- [ ] **Шаг 2:** `eventHandlers.ts` — по типу события патчить Query cache (`queryClient.setQueryData`)
- [ ] **Шаг 3:** Игнорировать `task.updated` если `payload.updatedAt <=` локальной версии (инициатор)
- [ ] **Шаг 4:** `BoardPage` — `useEffect` подключение при mount, отключение при unmount
- [ ] **Шаг 5:** UI: «Подключено» / «Переподключение…» / «N онлайн»

**Зачёт недели 7:** два браузера — изменение статуса видно без F5.

---

## Часть I. Неделя 8 — Ошибки и граничные случаи

### Task I1: Централизованная обработка ошибок

**Files:**
- Create: `frontend/src/shared/lib/errors.ts`, `frontend/src/shared/ui/Toast.tsx`
- Modify: `frontend/src/shared/api/client.ts`

- [ ] **Шаг 1:** Маппинг статусов

```typescript
export function handleApiError(error: { status: number; error?: string }) {
  if (error.status === 401) window.location.href = '/login';
  if (error.status === 409) toast('Циклическая зависимость запрещена');
  if (error.status === 403) toast('Нет доступа к проекту');
  // ...
}
```

- [ ] **Шаг 2:** ErrorBoundary вокруг `BoardCanvas`
- [ ] **Шаг 3:** Skeleton при `isLoading` на доске и списке проектов
- [ ] **Шаг 4:** Пустые состояния (нет проектов, нет задач)

---

## Часть J. Неделя 9 — Тесты и CI

### Task J1: Vitest config

**Files:**
- Modify: `frontend/vite.config.ts`
- Create: `.github/workflows/frontend-ci.yml`

- [ ] **Шаг 1:** Настроить vitest environment jsdom
- [ ] **Шаг 2:** Минимум 8 тестов (login, projects, cycleCheck×3, task form)
- [ ] **Шаг 3:** GitHub Actions

```yaml
name: frontend-ci
on: [pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: cd frontend && npm ci && npm run lint && npm test
```

- [ ] **Шаг 4:** README в `frontend/README.md` — установка, env, запуск

---

## Часть K. Неделя 10 — Развёртывание

### Task K1: Production build

- [ ] **Шаг 1:** `npm run build` — без ошибок
- [ ] **Шаг 2:** Deploy на Vercel/Netlify, env: `VITE_API_URL`, `VITE_WS_URL` (staging ментора)
- [ ] **Шаг 3:** Пройти чек-лист из спецификации §10
- [ ] **Шаг 4:** Записать demo 5 мин (скрипт в `docs/course/teams/{team}/demo-script.md`)
- [ ] **Шаг 5:** Ретроспектива — что получилось, что улучшить

---

## Чек-лист ментора по неделям

| Неделя | Готовность бекенда | Действие ментора |
|--------|-------------------|------------------|
| 1 | brief, OpenAPI draft, WS doc | ревью ТЗ |
| 2 | — | ревью макетов |
| 3 | auth на staging | помощь с CORS, JWT |
| 4 | projects + invites | проверить invite flow |
| 5 | GET tasks/deps | seed-данные для demo |
| 6 | CRUD + 409 cycle | пара тестовых сценариев цикла |
| 7 | WebSocket live | парное тестирование с командой |
| 8 | стабильный staging | ревью обработки ошибок |
| 9 | — | ревью CI и тестов |
| 10 | staging для demo | оценка по rubric из спецификации |

---

## Самопроверка плана

**1. Покрытие спецификации:**
- [x] Недели 1–10 — отдельные части B–K
- [x] Auth, проекты, invite — D, E
- [x] ReactFlow read/edit — F, G
- [x] WebSocket — H
- [x] Ошибки — I
- [x] Тесты, CI, deploy — J, K
- [x] Менторский бекенд — A1–A7

**2. Заглушки:** не обнаружены.

**3. Согласованность типов:** Task, Dependency, события WS совпадают между OpenAPI, mapToFlow и eventHandlers.

---

## Режим выполнения

**Выбрано:** только бекенд по неделям (часть A). Фронтенд (части B–K) выполняют студенты самостоятельно по плану; эталонный фронт в репозитории не создаётся.

**Порядок работы ментора:** Task A1 → A2 → A3 до старта курса; A4 на неделе 3; A5 на неделе 4; A6 на неделях 5–6; A7 на неделе 7.

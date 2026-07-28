-- name: ListTasksByProject :many
SELECT id, project_id, title, status, assignee_id, x, y, updated_at
FROM tasks
WHERE project_id = $1
ORDER BY created_at;

-- name: CreateTask :one
INSERT INTO tasks (project_id, title, status, assignee_id, x, y, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now())
RETURNING id, project_id, title, status, assignee_id, x, y, updated_at;

-- name: GetTask :one
SELECT id, project_id, title, status, assignee_id, x, y, updated_at
FROM tasks
WHERE id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET
    title = $2,
    status = $3,
    assignee_id = $4,
    x = $5,
    y = $6,
    updated_at = now()
WHERE id = $1
RETURNING id, project_id, title, status, assignee_id, x, y, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

-- name: GetTaskProjectID :one
SELECT project_id
FROM tasks
WHERE id = $1;

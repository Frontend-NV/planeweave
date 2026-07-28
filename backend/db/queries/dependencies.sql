-- name: ListDependenciesByProject :many
SELECT id, project_id, from_task_id, to_task_id
FROM dependencies
WHERE project_id = $1;

-- name: ListDependencyEdgesByProject :many
SELECT from_task_id, to_task_id
FROM dependencies
WHERE project_id = $1;

-- name: CreateDependency :one
INSERT INTO dependencies (project_id, from_task_id, to_task_id)
VALUES ($1, $2, $3)
RETURNING id, project_id, from_task_id, to_task_id;

-- name: GetDependencyProjectID :one
SELECT project_id
FROM dependencies
WHERE id = $1;

-- name: DeleteDependency :exec
DELETE FROM dependencies
WHERE id = $1;

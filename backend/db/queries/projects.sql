-- name: ListProjectsByUser :many
SELECT p.id, p.title, p.owner_id, p.created_at
FROM projects p
JOIN project_members pm ON pm.project_id = p.id
WHERE pm.user_id = $1
ORDER BY p.created_at DESC;

-- name: CreateProject :one
INSERT INTO projects (title, owner_id)
VALUES ($1, $2)
RETURNING id, title, owner_id, created_at;

-- name: GetProject :one
SELECT id, title, owner_id, created_at
FROM projects
WHERE id = $1;

-- name: UpdateProjectTitle :one
UPDATE projects
SET title = $1
WHERE id = $2
RETURNING id, title, owner_id, created_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: GetProjectOwnerID :one
SELECT owner_id
FROM projects
WHERE id = $1;

-- name: IsProjectMember :one
SELECT EXISTS(
    SELECT 1 FROM project_members
    WHERE project_id = $1 AND user_id = $2
) AS is_member;

-- name: AddProjectMember :exec
INSERT INTO project_members (project_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: ListProjectMembers :many
SELECT u.id, u.email, u.display_name, pm.role
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE pm.project_id = $1
ORDER BY pm.role DESC, u.display_name;

-- name: CreateInvite :exec
INSERT INTO invites (token, project_id, created_by, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetInvite :one
SELECT project_id, expires_at
FROM invites
WHERE token = $1;

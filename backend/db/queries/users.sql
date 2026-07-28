-- name: CreateUser :one
INSERT INTO users (email, password_hash, display_name)
VALUES ($1, $2, $3)
RETURNING id, email, display_name, created_at;

-- name: GetUserByEmail :one
SELECT id, email, display_name, password_hash
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, display_name
FROM users
WHERE id = $1;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: ListSessionsByUser :many
SELECT * FROM sessions
WHERE user_id = $1;

-- name: ListSessionsByModule :many
SELECT * FROM sessions
WHERE module_id = $1;

-- name: CreateSession :one
INSERT INTO sessions (tenant_id, module_id, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

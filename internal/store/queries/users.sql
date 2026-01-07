-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE tenant_id = $1 AND email = $2 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (tenant_id, email, display_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1;

-- name: UpdateUser :exec
UPDATE users
SET display_name = $2, is_active = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

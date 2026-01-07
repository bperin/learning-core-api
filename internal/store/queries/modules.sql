-- name: GetModule :one
SELECT * FROM modules
WHERE id = $1 LIMIT 1;

-- name: ListModulesByTenant :many
SELECT * FROM modules
WHERE tenant_id = $1;

-- name: GetModuleByName :one
SELECT * FROM modules
WHERE tenant_id = $1 AND name = $2 LIMIT 1;

-- name: CreateModule :one
INSERT INTO modules (tenant_id, title, name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateModule :exec
UPDATE modules SET name = $2, description = $3 WHERE id = $1;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1;

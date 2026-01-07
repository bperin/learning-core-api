-- name: GetTenantById :one
SELECT * FROM tenants
WHERE id = $1 LIMIT 1;

-- name: CreateTenant :one
INSERT INTO tenants (name, is_active)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllTenants :many
SELECT * FROM tenants;


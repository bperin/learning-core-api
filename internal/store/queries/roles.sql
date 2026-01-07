-- name: GetUserRoles :many
SELECT user_id, role, granted_at FROM user_roles
WHERE user_id = $1;

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role)
VALUES ($1, $2)
ON CONFLICT (user_id, role) DO NOTHING;

-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role = $2;

-- name: ListUsersByRole :many
SELECT u.* FROM users u
JOIN user_roles ur ON u.id = ur.user_id
WHERE u.tenant_id = $1 AND ur.role = $2;

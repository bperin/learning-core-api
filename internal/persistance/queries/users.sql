-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  id, email, password, is_admin, is_learner, is_teacher
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

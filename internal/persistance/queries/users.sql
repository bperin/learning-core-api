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

-- name: UpdateUser :one
UPDATE users SET
  email = COALESCE($2, email),
  password = COALESCE($3, password),
  is_admin = COALESCE($4, is_admin),
  is_learner = COALESCE($5, is_learner),
  is_teacher = COALESCE($6, is_teacher),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET
  password = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserRoles :one
UPDATE users SET
  is_admin = $2,
  is_learner = $3,
  is_teacher = $4,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE 
  (CASE WHEN $1::text = 'admin' THEN is_admin ELSE false END) OR
  (CASE WHEN $1::text = 'teacher' THEN is_teacher ELSE false END) OR
  (CASE WHEN $1::text = 'learner' THEN is_learner ELSE false END)
ORDER BY created_at DESC;

-- name: SearchUsersByEmail :many
SELECT * FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY email
LIMIT $2 OFFSET $3;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountUsersByRole :one
SELECT COUNT(*) FROM users
WHERE 
  (CASE WHEN $1::text = 'admin' THEN is_admin ELSE false END) OR
  (CASE WHEN $1::text = 'teacher' THEN is_teacher ELSE false END) OR
  (CASE WHEN $1::text = 'learner' THEN is_learner ELSE false END);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id, token, expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;

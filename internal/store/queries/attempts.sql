-- name: CreateAttempt :one
INSERT INTO attempts (session_id, tenant_id, artifact_id, is_correct, user_answer)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAttempt :one
SELECT * FROM attempts
WHERE id = $1 LIMIT 1;

-- name: ListAttemptsByTenant :many
SELECT * FROM attempts
WHERE tenant_id = $1;

-- name: ListAttemptsBySession :many
SELECT * FROM attempts
WHERE session_id = $1;

-- name: DeleteAttempt :exec
DELETE FROM attempts
WHERE id = $1;

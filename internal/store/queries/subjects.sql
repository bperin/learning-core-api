-- name: CreateSubject :one
INSERT INTO subjects (user_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSubject :one
SELECT * FROM subjects
WHERE id = $1 LIMIT 1;

-- name: GetSubjectByUserAndName :one
SELECT * FROM subjects
WHERE user_id = $1 AND name = $2 LIMIT 1;

-- name: ListSubjectsByUser :many
SELECT * FROM subjects
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSubject :one
UPDATE subjects
SET name = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSubject :exec
DELETE FROM subjects WHERE id = $1;

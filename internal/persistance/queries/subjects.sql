-- name: GetSubject :one
SELECT * FROM subjects
WHERE id = $1 LIMIT 1;

-- name: ListSubjectsByUser :many
SELECT * FROM subjects
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateSubject :one
INSERT INTO subjects (
  id, name, description, user_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateSubject :one
UPDATE subjects
SET name = $2,
    description = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSubject :exec
DELETE FROM subjects
WHERE id = $1;

-- name: ListAllSubjects :many
SELECT * FROM subjects
ORDER BY name;

-- name: SearchSubjectsByName :many
SELECT * FROM subjects
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CountSubjectsByUser :one
SELECT COUNT(*) FROM subjects
WHERE user_id = $1;

-- name: GetSubjectByUserAndName :one
SELECT * FROM subjects
WHERE user_id = $1 AND name = $2
LIMIT 1;

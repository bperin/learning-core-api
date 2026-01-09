-- name: GetCurriculum :one
SELECT * FROM curricula WHERE id = $1 LIMIT 1;

-- name: ListCurriculaBySubject :many
SELECT * FROM curricula
WHERE subject_id = $1
ORDER BY order_index ASC, created_at ASC;

-- name: ListCurriculaByParent :many
SELECT * FROM curricula
WHERE parent_id = $1
ORDER BY order_index ASC, created_at ASC;

-- name: ListRootCurriculaBySubject :many
SELECT * FROM curricula
WHERE subject_id = $1 AND parent_id IS NULL
ORDER BY order_index ASC, created_at ASC;

-- name: CreateCurriculum :one
INSERT INTO curricula (
  subject_id, parent_id, label, code, description, order_index,
  grade_level, is_active
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateCurriculum :one
UPDATE curricula SET
  subject_id = COALESCE($2, subject_id),
  parent_id = COALESCE($3, parent_id),
  label = COALESCE($4, label),
  code = COALESCE($5, code),
  description = COALESCE($6, description),
  order_index = COALESCE($7, order_index),
  grade_level = COALESCE($8, grade_level),
  is_active = COALESCE($9, is_active),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCurriculum :exec
DELETE FROM curricula WHERE id = $1;

-- name: GetEval :one
SELECT * FROM evals WHERE id = $1 LIMIT 1;

-- name: GetEvalsByUser :many
SELECT * FROM evals WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetEvalsBySubject :many
SELECT * FROM evals WHERE subject_id = $1 ORDER BY created_at DESC;

-- name: GetEvalsByStatus :many
SELECT * FROM evals WHERE status = $1 ORDER BY created_at DESC;

-- name: GetPublishedEvals :many
SELECT * FROM evals WHERE status = 'published' ORDER BY published_at DESC;

-- name: GetDraftEvals :many
SELECT * FROM evals WHERE status = 'draft' ORDER BY created_at DESC;

-- name: ListEvals :many
SELECT * FROM evals ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateEval :one
INSERT INTO evals (
  title, description, status, difficulty, instructions, rubric,
  subject_id, user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: PublishEval :one
UPDATE evals SET
  status = 'published',
  published_at = now(),
  updated_at = now()
WHERE id = $1 AND status = 'draft'
RETURNING *;

-- name: ArchiveEval :one
UPDATE evals SET
  status = 'archived',
  archived_at = now(),
  updated_at = now()
WHERE id = $1 AND status IN ('draft', 'published')
RETURNING *;

-- name: SearchEvalsByTitle :many
SELECT * FROM evals 
WHERE title ILIKE '%' || $1 || '%' 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetEvalWithItemCount :one
SELECT e.*, COUNT(ei.id) as item_count
FROM evals e
LEFT JOIN eval_items ei ON e.id = ei.eval_id
WHERE e.id = $1
GROUP BY e.id;

-- name: GetEvalsWithItemCounts :many
SELECT e.*, COUNT(ei.id) as item_count
FROM evals e
LEFT JOIN eval_items ei ON e.id = ei.eval_id
WHERE e.user_id = $1
GROUP BY e.id
ORDER BY e.created_at DESC;

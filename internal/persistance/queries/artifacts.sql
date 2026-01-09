-- name: GetArtifact :one
SELECT * FROM artifacts WHERE id = $1 LIMIT 1;

-- name: GetArtifactsByType :many
SELECT * FROM artifacts WHERE type = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByStatus :many
SELECT * FROM artifacts WHERE status = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByDocument :many
SELECT * FROM artifacts WHERE document_id = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByEval :many
SELECT * FROM artifacts WHERE eval_id = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByEvalItem :many
SELECT * FROM artifacts WHERE eval_item_id = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByAttempt :many
SELECT * FROM artifacts WHERE attempt_id = $1 ORDER BY created_at DESC;

-- name: GetArtifactsByUser :many
SELECT * FROM artifacts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListArtifacts :many
SELECT * FROM artifacts ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateArtifact :one
INSERT INTO artifacts (
  type, status, document_id, eval_id, eval_item_id, attempt_id, user_id,
  text, json, model, prompt, input_hash, meta, error
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: UpdateArtifact :one
UPDATE artifacts SET
  status = COALESCE($2, status),
  text = COALESCE($3, text),
  json = COALESCE($4, json),
  model = COALESCE($5, model),
  prompt = COALESCE($6, prompt),
  meta = COALESCE($7, meta),
  error = COALESCE($8, error)
WHERE id = $1
RETURNING *;

-- name: UpdateArtifactStatus :one
UPDATE artifacts SET
  status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteArtifact :exec
DELETE FROM artifacts WHERE id = $1;

-- name: GetArtifactsByTypeAndEntity :many
SELECT * FROM artifacts 
WHERE type = $1 
AND (
  (document_id = $2 AND $2 IS NOT NULL) OR
  (eval_id = $3 AND $3 IS NOT NULL) OR
  (eval_item_id = $4 AND $4 IS NOT NULL) OR
  (attempt_id = $5 AND $5 IS NOT NULL)
)
ORDER BY created_at DESC;

-- name: GetArtifactsByInputHash :many
SELECT * FROM artifacts 
WHERE input_hash = $1 
ORDER BY created_at DESC;

-- name: GetLatestArtifactByTypeAndEntity :one
SELECT * FROM artifacts 
WHERE type = $1 
AND (
  (document_id = $2 AND $2 IS NOT NULL) OR
  (eval_id = $3 AND $3 IS NOT NULL) OR
  (eval_item_id = $4 AND $4 IS NOT NULL) OR
  (attempt_id = $5 AND $5 IS NOT NULL)
)
ORDER BY created_at DESC 
LIMIT 1;

-- name: GetArtifactStats :one
SELECT 
  COUNT(*) as total_artifacts,
  COUNT(CASE WHEN status = 'READY' THEN 1 END) as ready_count,
  COUNT(CASE WHEN status = 'PENDING' THEN 1 END) as pending_count,
  COUNT(CASE WHEN status = 'ERROR' THEN 1 END) as error_count,
  COUNT(CASE WHEN error IS NOT NULL THEN 1 END) as with_errors
FROM artifacts;

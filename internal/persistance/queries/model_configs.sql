-- name: CreateModelConfig :one
INSERT INTO model_configs (
  version, model_name, temperature, max_tokens, top_p, top_k, mime_type, is_active, created_by
) VALUES (
  (SELECT COALESCE(MAX(version), 0) + 1 FROM model_configs),
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetActiveModelConfig :one
SELECT * FROM model_configs WHERE is_active = true LIMIT 1;

-- name: ActivateModelConfig :exec
UPDATE model_configs SET is_active = true WHERE id = $1;

-- name: DeactivateOtherModelConfigs :exec
UPDATE model_configs SET is_active = false WHERE id != $1;

-- name: GetModelConfig :one
SELECT * FROM model_configs WHERE id = $1 LIMIT 1;

-- name: ListModelConfigs :many
SELECT * FROM model_configs ORDER BY created_at DESC;

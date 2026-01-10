-- name: CreateModelConfig :one
WITH inserted AS (
  INSERT INTO model_configs (
    version, model_name, temperature, max_tokens, top_p, top_k, mime_type, is_active, created_by
  ) VALUES (
    (SELECT COALESCE(MAX(version), 0) + 1 FROM model_configs),
    $1, $2, $3, $4, $5, $6, $7, $8
  )
  RETURNING *
),
deactivated AS (
  UPDATE model_configs SET
    is_active = false
  WHERE model_configs.id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: GetActiveModelConfig :one
SELECT * FROM model_configs WHERE is_active = true LIMIT 1;

-- name: ActivateModelConfig :exec
WITH activated AS (
  UPDATE model_configs SET is_active = true WHERE model_configs.id = $1
)
UPDATE model_configs SET is_active = false WHERE model_configs.id != $1;

-- name: DeactivateOtherModelConfigs :exec
UPDATE model_configs SET is_active = false WHERE id != $1;

-- name: GetModelConfig :one
SELECT * FROM model_configs WHERE id = $1 LIMIT 1;

-- name: ListModelConfigs :many
SELECT * FROM model_configs ORDER BY created_at DESC;

-- name: CreateChunkingConfig :one
WITH inserted AS (
  INSERT INTO chunking_configs (
    version, chunk_size, chunk_overlap, is_active, created_by
  ) VALUES (
    (SELECT COALESCE(MAX(version), 0) + 1 FROM chunking_configs),
    $1, $2, $3, $4
  )
  RETURNING *
),
deactivated AS (
  UPDATE chunking_configs SET
    is_active = false
  WHERE chunking_configs.id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: GetActiveChunkingConfig :one
SELECT * FROM chunking_configs WHERE is_active = true LIMIT 1;

-- name: ActivateChunkingConfig :exec
WITH activated AS (
  UPDATE chunking_configs SET is_active = true WHERE chunking_configs.id = $1
)
UPDATE chunking_configs SET is_active = false WHERE chunking_configs.id != $1;

-- name: DeactivateOtherChunkingConfigs :exec
UPDATE chunking_configs SET is_active = false WHERE id != $1;

-- name: GetChunkingConfig :one
SELECT * FROM chunking_configs WHERE id = $1 LIMIT 1;

-- name: ListChunkingConfigs :many
SELECT * FROM chunking_configs ORDER BY created_at DESC;

-- name: CreateChunkingConfig :one
INSERT INTO chunking_configs (
  version, chunk_size, chunk_overlap, is_active, created_by
) VALUES (
  (SELECT COALESCE(MAX(version), 0) + 1 FROM chunking_configs),
  $1, $2, $3, $4
) RETURNING *;

-- name: GetActiveChunkingConfig :one
SELECT * FROM chunking_configs WHERE is_active = true LIMIT 1;

-- name: ActivateChunkingConfig :exec
UPDATE chunking_configs SET is_active = true WHERE id = $1;

-- name: DeactivateOtherChunkingConfigs :exec
UPDATE chunking_configs SET is_active = false WHERE id != $1;

-- name: GetChunkingConfig :one
SELECT * FROM chunking_configs WHERE id = $1 LIMIT 1;

-- name: ListChunkingConfigs :many
SELECT * FROM chunking_configs ORDER BY created_at DESC;

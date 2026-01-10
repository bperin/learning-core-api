-- name: CreateSystemInstruction :one
INSERT INTO system_instructions (
  version, text, is_active, created_by
) VALUES (
  (SELECT COALESCE(MAX(version), 0) + 1 FROM system_instructions),
  $1, $2, $3
) RETURNING *;

-- name: GetActiveSystemInstruction :one
SELECT * FROM system_instructions WHERE is_active = true LIMIT 1;

-- name: ActivateSystemInstruction :exec
UPDATE system_instructions SET is_active = true WHERE id = $1;

-- name: DeactivateOtherSystemInstructions :exec
UPDATE system_instructions SET is_active = false WHERE id != $1;

-- name: GetSystemInstruction :one
SELECT * FROM system_instructions WHERE id = $1 LIMIT 1;

-- name: ListSystemInstructions :many
SELECT * FROM system_instructions ORDER BY created_at DESC;

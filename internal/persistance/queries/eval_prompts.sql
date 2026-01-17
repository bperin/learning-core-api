-- name: GetEvalPrompt :one
SELECT * FROM eval_prompts WHERE id = $1 LIMIT 1;

-- name: GetActiveEvalPrompt :one
SELECT * FROM eval_prompts WHERE eval_type = $1 AND is_active = true ORDER BY version DESC LIMIT 1;

-- name: GetEvalPromptByVersion :one
SELECT * FROM eval_prompts WHERE eval_type = $1 AND version = $2 LIMIT 1;

-- name: ListEvalPrompts :many
SELECT * FROM eval_prompts WHERE eval_type = $1 ORDER BY version DESC LIMIT $2 OFFSET $3;

-- name: CreateEvalPrompt :one
INSERT INTO eval_prompts (
  eval_type, version, prompt_text, description, is_active, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: DeactivateEvalPrompt :exec
UPDATE eval_prompts SET is_active = false WHERE id = $1;

-- name: ActivateEvalPrompt :exec
UPDATE eval_prompts SET is_active = true WHERE id = $1;

-- name: GetLatestEvalPromptVersion :one
SELECT COALESCE(MAX(version), 0) as latest_version FROM eval_prompts WHERE eval_type = $1;

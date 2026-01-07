-- name: CreateGenerationRun :one
INSERT INTO generation_runs (module_id, agent_name, agent_version, model, model_params, prompt_id, store_name, metadata_filter, status, input_payload)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetGenerationRun :one
SELECT * FROM generation_runs
WHERE id = $1 LIMIT 1;

-- name: ListGenerationRunsByModule :many
SELECT * FROM generation_runs
WHERE module_id = $1
ORDER BY created_at DESC;

-- name: ListGenerationRunsByStatus :many
SELECT * FROM generation_runs
WHERE module_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: UpdateGenerationRun :exec
UPDATE generation_runs
SET status = $2, output_payload = $3, error = $4, started_at = $5, finished_at = $6
WHERE id = $1;

-- name: DeleteGenerationRun :exec
DELETE FROM generation_runs WHERE id = $1;

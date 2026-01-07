-- name: CreateEvalRun :one
INSERT INTO eval_runs (
  artifact_id,
  generation_run_id,
  suite_id,
  judge_model,
  judge_params,
  status,
  started_at
) VALUES (
  $1,        -- artifact_id
  $2,        -- generation_run_id (nullable)
  $3,        -- suite_id
  $4,        -- judge_model
  $5,        -- judge_params :: jsonb
  'RUNNING',
  now()
)
RETURNING *;

-- name: GetEvalRunByID :one
SELECT *
FROM eval_runs
WHERE id = $1
LIMIT 1;

-- name: GetLatestEvalRunForArtifact :one
SELECT *
FROM eval_runs
WHERE artifact_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListEvalRunsByArtifact :many
SELECT *
FROM eval_runs
WHERE artifact_id = $1
ORDER BY created_at DESC;

-- name: UpdateEvalRunResult :one
UPDATE eval_runs
SET
  status = $2,          -- run_status
  overall_pass = $3,    -- boolean
  overall_score = $4,   -- real
  finished_at = now(),
  error = $5            -- jsonb (nullable)
WHERE id = $1
RETURNING *;

-- name: DeleteEvalRun :exec
DELETE FROM eval_runs
WHERE id = $1;

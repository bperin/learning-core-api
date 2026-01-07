-- name: UpsertEvalResult :one
INSERT INTO eval_results (
  eval_run_id,
  rule_id,
  pass,
  score,
  details
) VALUES (
  $1,        -- eval_run_id
  $2,        -- rule_id
  $3,        -- pass
  $4,        -- score
  $5         -- details :: jsonb
)
ON CONFLICT (eval_run_id, rule_id)
DO UPDATE SET
  pass = EXCLUDED.pass,
  score = EXCLUDED.score,
  details = EXCLUDED.details
RETURNING *;

-- name: GetEvalResultByID :one
SELECT *
FROM eval_results
WHERE id = $1
LIMIT 1;

-- name: ListEvalResultsByRun :many
SELECT *
FROM eval_results
WHERE eval_run_id = $1
ORDER BY created_at ASC;

-- name: ListEvalResultsByRule :many
SELECT *
FROM eval_results
WHERE rule_id = $1
ORDER BY created_at DESC;

-- name: DeleteEvalResultsByRun :exec
DELETE FROM eval_results
WHERE eval_run_id = $1;

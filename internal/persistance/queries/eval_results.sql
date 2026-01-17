-- name: GetEvalResult :one
SELECT * FROM eval_results WHERE id = $1 LIMIT 1;

-- name: GetEvalResultsByEvalItem :many
SELECT * FROM eval_results WHERE eval_item_id = $1 ORDER BY created_at DESC;

-- name: GetEvalResultsByType :many
SELECT * FROM eval_results WHERE eval_type = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListEvalResults :many
SELECT * FROM eval_results ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateEvalResult :one
INSERT INTO eval_results (
  eval_item_id, eval_type, eval_prompt_id, score, is_grounded, verdict, reasoning, unsupported_claims, gcp_eval_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetLatestEvalResultForItem :one
SELECT * FROM eval_results WHERE eval_item_id = $1 AND eval_type = $2 ORDER BY created_at DESC LIMIT 1;

-- name: GetEvalResultStats :one
SELECT 
  COUNT(*) as total_evals,
  COUNT(CASE WHEN verdict = 'PASS' THEN 1 END) as passed,
  COUNT(CASE WHEN verdict = 'FAIL' THEN 1 END) as failed,
  COUNT(CASE WHEN verdict = 'WARN' THEN 1 END) as warned,
  ROUND(AVG(score)::numeric, 2) as avg_score
FROM eval_results 
WHERE eval_type = $1;

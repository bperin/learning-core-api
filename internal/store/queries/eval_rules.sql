-- name: CreateEvalRule :one
INSERT INTO eval_rules (
  suite_id,
  eval_type,
  min_score,
  max_score,
  weight,
  hard_fail,
  params
) VALUES (
  $1,        -- suite_id
  $2,        -- eval_type
  $3,        -- min_score
  $4,        -- max_score
  $5,        -- weight
  $6,        -- hard_fail
  $7         -- params :: jsonb
)
RETURNING *;

-- name: GetEvalRule :one
SELECT *
FROM eval_rules
WHERE id = $1
LIMIT 1;

-- name: GetEvalRuleBySuiteAndType :one
SELECT *
FROM eval_rules
WHERE suite_id = $1
  AND eval_type = $2
LIMIT 1;

-- name: ListEvalRulesBySuite :many
SELECT *
FROM eval_rules
WHERE suite_id = $1
ORDER BY created_at ASC;

-- name: ListEvalRulesByEvalType :many
SELECT *
FROM eval_rules
WHERE eval_type = $1
ORDER BY created_at DESC;

-- name: UpdateEvalRule :one
UPDATE eval_rules
SET
  min_score = $2,
  max_score = $3,
  weight = $4,
  hard_fail = $5,
  params = $6
WHERE id = $1
RETURNING *;

-- name: DeleteEvalRule :exec
DELETE FROM eval_rules
WHERE id = $1;

-- name: DeleteEvalRulesBySuite :exec
DELETE FROM eval_rules
WHERE suite_id = $1;

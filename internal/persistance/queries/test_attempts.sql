-- name: GetTestAttempt :one
SELECT * FROM test_attempts WHERE id = $1 LIMIT 1;

-- name: GetTestAttemptsByUser :many
SELECT * FROM test_attempts WHERE user_id = $1 ORDER BY started_at DESC;

-- name: GetTestAttemptsByEval :many
SELECT * FROM test_attempts WHERE eval_id = $1 ORDER BY started_at DESC;

-- name: GetCompletedAttempts :many
SELECT * FROM test_attempts WHERE completed_at IS NOT NULL ORDER BY completed_at DESC;

-- name: GetActiveAttempts :many
SELECT * FROM test_attempts WHERE completed_at IS NULL ORDER BY started_at DESC;

-- name: GetUserAttemptsByEval :many
SELECT * FROM test_attempts 
WHERE user_id = $1 AND eval_id = $2 
ORDER BY started_at DESC;

-- name: ListTestAttempts :many
SELECT * FROM test_attempts ORDER BY started_at DESC LIMIT $1 OFFSET $2;

-- name: CreateTestAttempt :one
INSERT INTO test_attempts (
  user_id, eval_id, total
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateTestAttemptScore :one
UPDATE test_attempts SET
  score = $2,
  percentage = CASE 
    WHEN total > 0 THEN ROUND(($2::numeric / total::numeric) * 100, 2)
    ELSE 0 
  END,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CompleteTestAttempt :one
UPDATE test_attempts SET
  score = $2,
  percentage = CASE 
    WHEN total > 0 THEN ROUND(($2::numeric / total::numeric) * 100, 2)
    ELSE 0 
  END,
  total_time = $3,
  feedback = $4,
  summary = $5,
  completed_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateTestAttemptTime :one
UPDATE test_attempts SET
  total_time = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTestAttempt :exec
DELETE FROM test_attempts WHERE id = $1;

-- name: GetTestAttemptWithAnswers :one
SELECT 
  ta.*,
  COUNT(ua.id) as answer_count,
  COUNT(CASE WHEN ua.is_correct = true THEN 1 END) as correct_count
FROM test_attempts ta
LEFT JOIN user_answers ua ON ta.id = ua.attempt_id
WHERE ta.id = $1
GROUP BY ta.id;

-- name: GetUserTestStats :one
SELECT 
  COUNT(*) as total_attempts,
  COUNT(CASE WHEN completed_at IS NOT NULL THEN 1 END) as completed_attempts,
  AVG(CASE WHEN percentage IS NOT NULL THEN percentage END) as avg_percentage,
  MAX(percentage) as best_percentage,
  AVG(CASE WHEN total_time IS NOT NULL THEN total_time END) as avg_time
FROM test_attempts 
WHERE user_id = $1;

-- name: GetEvalTestStats :one
SELECT 
  COUNT(*) as total_attempts,
  COUNT(CASE WHEN completed_at IS NOT NULL THEN 1 END) as completed_attempts,
  AVG(CASE WHEN percentage IS NOT NULL THEN percentage END) as avg_percentage,
  COUNT(DISTINCT user_id) as unique_users
FROM test_attempts 
WHERE eval_id = $1;

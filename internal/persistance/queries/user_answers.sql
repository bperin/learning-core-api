-- name: GetUserAnswer :one
SELECT * FROM user_answers WHERE id = $1 LIMIT 1;

-- name: GetUserAnswersByAttempt :many
SELECT * FROM user_answers WHERE attempt_id = $1 ORDER BY created_at ASC;

-- name: GetUserAnswersByEvalItem :many
SELECT * FROM user_answers WHERE eval_item_id = $1 ORDER BY created_at DESC;

-- name: GetUserAnswerByAttemptAndItem :one
SELECT * FROM user_answers 
WHERE attempt_id = $1 AND eval_item_id = $2 
LIMIT 1;

-- name: ListUserAnswers :many
SELECT * FROM user_answers ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateUserAnswer :one
INSERT INTO user_answers (
  attempt_id, eval_item_id, selected_idx, is_correct, time_spent, hints_used
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAnswersByUserAndEval :many
SELECT ua.*, ei.prompt, ei.options, ei.correct_idx, ei.explanation
FROM user_answers ua
JOIN test_attempts ta ON ua.attempt_id = ta.id
JOIN eval_items ei ON ua.eval_item_id = ei.id
WHERE ta.user_id = $1 AND ta.eval_id = $2
ORDER BY ua.created_at ASC;

-- name: GetCorrectAnswersByAttempt :many
SELECT * FROM user_answers 
WHERE attempt_id = $1 AND is_correct = true 
ORDER BY created_at ASC;

-- name: GetIncorrectAnswersByAttempt :many
SELECT * FROM user_answers 
WHERE attempt_id = $1 AND is_correct = false 
ORDER BY created_at ASC;

-- name: GetAnswerStatsForEvalItem :one
SELECT 
  COUNT(*) as total_answers,
  COUNT(CASE WHEN is_correct = true THEN 1 END) as correct_answers,
  COUNT(CASE WHEN is_correct = false THEN 1 END) as incorrect_answers,
  CASE 
    WHEN COUNT(*) > 0 THEN 
      ROUND(COUNT(CASE WHEN is_correct = true THEN 1 END)::numeric / COUNT(*)::numeric * 100, 2)
    ELSE 0 
  END as success_rate,
  AVG(time_spent) as avg_time_spent,
  AVG(hints_used) as avg_hints_used
FROM user_answers 
WHERE eval_item_id = $1;

-- name: GetUserAnswerPatterns :many
SELECT 
  selected_idx,
  COUNT(*) as selection_count,
  COUNT(CASE WHEN is_correct = true THEN 1 END) as correct_count
FROM user_answers 
WHERE eval_item_id = $1
GROUP BY selected_idx
ORDER BY selected_idx;

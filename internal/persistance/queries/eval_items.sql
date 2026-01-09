-- name: GetEvalItem :one
SELECT * FROM eval_items WHERE id = $1 LIMIT 1;

-- name: GetEvalItemsByEval :many
SELECT * FROM eval_items WHERE eval_id = $1 ORDER BY id ASC;

-- name: ListEvalItems :many
SELECT * FROM eval_items ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: CreateEvalItem :one
INSERT INTO eval_items (
  eval_id, prompt, options, correct_idx, hint, explanation, metadata
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateEvalItem :one
UPDATE eval_items SET
  prompt = COALESCE($2, prompt),
  options = COALESCE($3, options),
  correct_idx = COALESCE($4, correct_idx),
  hint = COALESCE($5, hint),
  explanation = COALESCE($6, explanation),
  metadata = COALESCE($7, metadata),
  updated_at = now()
WHERE eval_items.id = $1
AND eval_items.eval_id IN (SELECT e.id FROM evals e WHERE e.status = 'draft')
RETURNING *;

-- name: DeleteEvalItem :exec
DELETE FROM eval_items 
WHERE eval_items.id = $1 
AND eval_items.eval_id IN (SELECT e.id FROM evals e WHERE e.status = 'draft');

-- name: GetEvalItemsWithAnswerStats :many
SELECT 
  ei.*,
  COUNT(ua.id) as total_answers,
  COUNT(CASE WHEN ua.is_correct = true THEN 1 END) as correct_answers,
  CASE 
    WHEN COUNT(ua.id) > 0 THEN 
      ROUND(COUNT(CASE WHEN ua.is_correct = true THEN 1 END)::numeric / COUNT(ua.id)::numeric * 100, 2)
    ELSE 0 
  END as success_rate
FROM eval_items ei
LEFT JOIN user_answers ua ON ei.id = ua.eval_item_id
WHERE ei.eval_id = $1
GROUP BY ei.id
ORDER BY ei.id ASC;

-- name: GetEvalItemWithReviews :one
SELECT 
  ei.*,
  COUNT(eir.id) as review_count,
  COUNT(CASE WHEN eir.verdict = 'APPROVED' THEN 1 END) as approved_count,
  COUNT(CASE WHEN eir.verdict = 'REJECTED' THEN 1 END) as rejected_count,
  COUNT(CASE WHEN eir.verdict = 'NEEDS_REVISION' THEN 1 END) as needs_revision_count
FROM eval_items ei
LEFT JOIN eval_item_reviews eir ON ei.id = eir.eval_item_id
WHERE ei.id = $1
GROUP BY ei.id;

-- name: SearchEvalItemsByPrompt :many
SELECT * FROM eval_items 
WHERE prompt ILIKE '%' || $1 || '%' 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetRandomEvalItems :many
SELECT * FROM eval_items 
WHERE eval_id = $1 
ORDER BY RANDOM() 
LIMIT $2;

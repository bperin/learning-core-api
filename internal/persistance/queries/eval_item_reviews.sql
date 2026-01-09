-- name: GetEvalItemReview :one
SELECT * FROM eval_item_reviews WHERE id = $1 LIMIT 1;

-- name: GetReviewsByEvalItem :many
SELECT * FROM eval_item_reviews WHERE eval_item_id = $1 ORDER BY created_at DESC;

-- name: GetReviewsByReviewer :many
SELECT * FROM eval_item_reviews WHERE reviewer_id = $1 ORDER BY created_at DESC;

-- name: GetReviewsByVerdict :many
SELECT * FROM eval_item_reviews WHERE verdict = $1 ORDER BY created_at DESC;

-- name: ListEvalItemReviews :many
SELECT * FROM eval_item_reviews ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateEvalItemReview :one
INSERT INTO eval_item_reviews (
  eval_item_id, reviewer_id, verdict, reasons, comments
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetReviewsByEvalItemAndReviewer :one
SELECT * FROM eval_item_reviews 
WHERE eval_item_id = $1 AND reviewer_id = $2 
ORDER BY created_at DESC 
LIMIT 1;

-- name: GetReviewStatsForEvalItem :one
SELECT 
  COUNT(*) as total_reviews,
  COUNT(CASE WHEN verdict = 'APPROVED' THEN 1 END) as approved_count,
  COUNT(CASE WHEN verdict = 'REJECTED' THEN 1 END) as rejected_count,
  COUNT(CASE WHEN verdict = 'NEEDS_REVISION' THEN 1 END) as needs_revision_count,
  CASE 
    WHEN COUNT(*) > 0 THEN 
      ROUND(COUNT(CASE WHEN verdict = 'APPROVED' THEN 1 END)::numeric / COUNT(*)::numeric * 100, 2)
    ELSE 0 
  END as approval_rate
FROM eval_item_reviews 
WHERE eval_item_id = $1;

-- name: GetReviewStatsForReviewer :one
SELECT 
  COUNT(*) as total_reviews,
  COUNT(CASE WHEN verdict = 'APPROVED' THEN 1 END) as approved_count,
  COUNT(CASE WHEN verdict = 'REJECTED' THEN 1 END) as rejected_count,
  COUNT(CASE WHEN verdict = 'NEEDS_REVISION' THEN 1 END) as needs_revision_count
FROM eval_item_reviews 
WHERE reviewer_id = $1;

-- name: GetReviewsWithEvalItemDetails :many
SELECT 
  eir.*,
  ei.prompt,
  ei.options,
  ei.correct_idx,
  e.title as eval_title
FROM eval_item_reviews eir
JOIN eval_items ei ON eir.eval_item_id = ei.id
JOIN evals e ON ei.eval_id = e.id
WHERE eir.reviewer_id = $1
ORDER BY eir.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPendingReviewsForEval :many
SELECT DISTINCT ei.*
FROM eval_items ei
LEFT JOIN eval_item_reviews eir ON ei.id = eir.eval_item_id
WHERE ei.eval_id = $1
AND (eir.id IS NULL OR eir.verdict = 'NEEDS_REVISION')
ORDER BY ei.id ASC;

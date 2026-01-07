-- name: CreateEvalSuite :one
INSERT INTO eval_suites (
  name,
  description
) VALUES (
  $1,        -- name
  $2         -- description
)
RETURNING *;

-- name: GetEvalSuiteByID :one
SELECT *
FROM eval_suites
WHERE id = $1
LIMIT 1;

-- name: GetEvalSuiteByName :one
SELECT *
FROM eval_suites
WHERE name = $1
LIMIT 1;

-- name: ListEvalSuites :many
SELECT *
FROM eval_suites
ORDER BY created_at DESC;

-- name: UpdateEvalSuite :one
UPDATE eval_suites
SET
  name = $2,
  description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEvalSuite :exec
DELETE FROM eval_suites
WHERE id = $1;

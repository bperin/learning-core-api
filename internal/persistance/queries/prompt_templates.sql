-- name: GetPromptTemplate :one
SELECT * FROM prompt_templates WHERE id = $1 LIMIT 1;

-- name: GetPromptTemplateByKey :one
SELECT * FROM prompt_templates WHERE key = $1 AND is_active = true LIMIT 1;

-- name: GetPromptTemplateByKeyAndVersion :one
SELECT * FROM prompt_templates WHERE key = $1 AND version = $2 LIMIT 1;

-- name: GetPromptTemplatesByKey :many
SELECT * FROM prompt_templates WHERE key = $1 ORDER BY version DESC;

-- name: GetActivePromptTemplates :many
SELECT * FROM prompt_templates WHERE is_active = true ORDER BY key ASC;

-- name: ListPromptTemplates :many
SELECT * FROM prompt_templates ORDER BY key ASC, version DESC LIMIT $1 OFFSET $2;

-- name: CreatePromptTemplate :one
INSERT INTO prompt_templates (
  key, version, is_active, title, description, template, metadata, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ActivatePromptTemplate :one
UPDATE prompt_templates SET
  is_active = true,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivatePromptTemplate :one
UPDATE prompt_templates SET
  is_active = false,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeactivateOtherVersions :exec
UPDATE prompt_templates SET
  is_active = false,
  updated_at = now()
WHERE key = $1 AND id != $2;

-- name: GetLatestVersionByKey :one
SELECT COALESCE(MAX(version), 0) as latest_version
FROM prompt_templates 
WHERE key = $1;

-- name: CreateNewVersion :one
INSERT INTO prompt_templates (
  key, version, is_active, title, description, template, metadata, created_by
) VALUES (
  $1, 
  (SELECT COALESCE(MAX(version), 0) + 1 FROM prompt_templates WHERE key = $1),
  $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: SearchPromptTemplatesByTitle :many
SELECT * FROM prompt_templates 
WHERE title ILIKE '%' || $1 || '%' 
ORDER BY key ASC, version DESC 
LIMIT $2 OFFSET $3;

-- name: GetPromptTemplatesByCreator :many
SELECT * FROM prompt_templates 
WHERE created_by = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetPromptTemplateStats :one
SELECT 
  COUNT(*) as total_templates,
  COUNT(DISTINCT key) as unique_keys,
  COUNT(CASE WHEN is_active = true THEN 1 END) as active_templates,
  AVG(version) as avg_version
FROM prompt_templates;

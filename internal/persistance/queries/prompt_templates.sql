-- name: GetPromptTemplate :one
SELECT * FROM prompt_templates WHERE id = $1 LIMIT 1;

-- name: GetPromptTemplateByGenerationType :one
SELECT * FROM prompt_templates WHERE generation_type = $1 AND is_active = true LIMIT 1;

-- name: GetPromptTemplateByGenerationTypeAndVersion :one
SELECT * FROM prompt_templates WHERE generation_type = $1 AND version = $2 LIMIT 1;

-- name: GetPromptTemplatesByGenerationType :many
SELECT * FROM prompt_templates WHERE generation_type = $1 ORDER BY version DESC;

-- name: GetActivePromptTemplates :many
SELECT * FROM prompt_templates WHERE is_active = true ORDER BY generation_type ASC;

-- name: ListPromptTemplates :many
SELECT * FROM prompt_templates ORDER BY generation_type ASC, version DESC LIMIT $1 OFFSET $2;

-- name: CreatePromptTemplate :one
WITH inserted AS (
  INSERT INTO prompt_templates (
    generation_type, version, is_active, title, description, template, metadata, created_by
  ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
  )
  RETURNING *
),
deactivated AS (
  UPDATE prompt_templates SET
    is_active = false,
    updated_at = now()
  WHERE generation_type = (SELECT generation_type FROM inserted)
    AND id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: ActivatePromptTemplate :one
WITH target AS (
  SELECT generation_type FROM prompt_templates WHERE prompt_templates.id = $1
),
deactivated AS (
  UPDATE prompt_templates SET
    is_active = false,
    updated_at = now()
  WHERE generation_type = (SELECT generation_type FROM target) AND id != $1
),
activated AS (
  UPDATE prompt_templates SET
    is_active = true,
    updated_at = now()
  WHERE id = $1
  RETURNING *
)
SELECT * FROM activated;

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
WHERE generation_type = $1 AND id != $2;

-- name: GetLatestVersionByGenerationType :one
SELECT COALESCE(MAX(version), 0) as latest_version
FROM prompt_templates 
WHERE generation_type = $1;

-- name: CreateNewVersion :one
WITH inserted AS (
  INSERT INTO prompt_templates (
    generation_type, version, is_active, title, description, template, metadata, created_by
  ) VALUES (
    $1, 
    (SELECT COALESCE(MAX(version), 0) + 1 FROM prompt_templates WHERE generation_type = $1),
    $2, $3, $4, $5, $6, $7
  )
  RETURNING *
),
deactivated AS (
  UPDATE prompt_templates SET
    is_active = false,
    updated_at = now()
  WHERE generation_type = (SELECT generation_type FROM inserted)
    AND id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: SearchPromptTemplatesByTitle :many
SELECT * FROM prompt_templates 
WHERE title ILIKE '%' || $1 || '%' 
ORDER BY generation_type ASC, version DESC 
LIMIT $2 OFFSET $3;

-- name: GetPromptTemplatesByCreator :many
SELECT * FROM prompt_templates 
WHERE created_by = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetPromptTemplateStats :one
SELECT 
  COUNT(*) as total_templates,
  COUNT(DISTINCT generation_type) as unique_keys,
  COUNT(CASE WHEN is_active = true THEN 1 END) as active_templates,
  AVG(version) as avg_version
FROM prompt_templates;

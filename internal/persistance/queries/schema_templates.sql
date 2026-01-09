-- name: GetSchemaTemplate :one
SELECT * FROM schema_templates WHERE id = $1 LIMIT 1;

-- name: GetSchemaTemplateByTypeAndVersion :one
SELECT * FROM schema_templates WHERE schema_type = $1 AND version = $2 LIMIT 1;

-- name: GetActiveSchemaTemplateByType :one
SELECT * FROM schema_templates WHERE schema_type = $1 AND is_active = true LIMIT 1;

-- name: ListSchemaTemplatesByType :many
SELECT * FROM schema_templates
WHERE schema_type = $1
ORDER BY version DESC;

-- name: ListActiveSchemaTemplates :many
SELECT * FROM schema_templates
WHERE is_active = true
ORDER BY schema_type ASC, version DESC;

-- name: ListSchemaTemplatesBySubject :many
SELECT * FROM schema_templates
WHERE subject_id = $1
ORDER BY schema_type ASC, version DESC;

-- name: ListSchemaTemplatesByCurriculum :many
SELECT * FROM schema_templates
WHERE curriculum_id = $1
ORDER BY schema_type ASC, version DESC;

-- name: CreateSchemaTemplate :one
INSERT INTO schema_templates (
  schema_type, version, schema_json, subject_id, curriculum_id,
  is_active, created_by, locked_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: CreateSchemaTemplateVersion :one
INSERT INTO schema_templates (
  schema_type, version, schema_json, subject_id, curriculum_id,
  is_active, created_by, locked_at
) VALUES (
  $1,
  (SELECT COALESCE(MAX(version), 0) + 1 FROM schema_templates WHERE schema_type = $1),
  $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ActivateSchemaTemplate :one
UPDATE schema_templates SET
  is_active = true
WHERE id = $1
RETURNING *;

-- name: DeactivateSchemaTemplate :one
UPDATE schema_templates SET
  is_active = false,
  locked_at = COALESCE(locked_at, now())
WHERE id = $1
RETURNING *;

-- name: DeactivateOtherSchemaTemplateVersions :exec
UPDATE schema_templates SET
  is_active = false,
  locked_at = COALESCE(locked_at, now())
WHERE schema_type = $1 AND id != $2;

-- name: GetSchemaTemplate :one
SELECT * FROM schema_templates WHERE id = $1 LIMIT 1;

-- name: GetActiveSchemaTemplateByGenerationType :one
SELECT * FROM schema_templates WHERE generation_type = $1 AND is_active = true LIMIT 1;

-- name: ListSchemaTemplatesByGenerationType :many
SELECT * FROM schema_templates
WHERE generation_type = $1
ORDER BY version DESC;

-- name: ListActiveSchemaTemplates :many
SELECT * FROM schema_templates
WHERE is_active = true
ORDER BY generation_type ASC, version DESC;

-- name: CreateSchemaTemplate :one
WITH inserted AS (
  INSERT INTO schema_templates (
    generation_type, version, schema_json,
    is_active, created_by, locked_at
  ) VALUES (
    $1,
    (SELECT COALESCE(MAX(version), 0) + 1 FROM schema_templates WHERE generation_type = $1),
    $2, $3, $4, $5
  )
  RETURNING *
),
deactivated AS (
  UPDATE schema_templates SET
    is_active = false,
    locked_at = COALESCE(locked_at, now())
  WHERE generation_type = (SELECT generation_type FROM inserted)
    AND id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: ActivateSchemaTemplate :one
WITH target AS (
  SELECT generation_type FROM schema_templates WHERE schema_templates.id = $1
),
deactivated AS (
  UPDATE schema_templates SET
    is_active = false,
    locked_at = COALESCE(locked_at, now())
  WHERE generation_type = (SELECT generation_type FROM target) AND id != $1
),
activated AS (
  UPDATE schema_templates SET
    is_active = true
  WHERE id = $1
  RETURNING *
)
SELECT * FROM activated;

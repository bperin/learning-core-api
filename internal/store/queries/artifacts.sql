-- name: CreateArtifact :one
INSERT INTO artifacts (module_id, generation_run_id, type, status, schema_version, difficulty, tags, artifact_payload, grounding, evidence_version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetArtifact :one
SELECT * FROM artifacts
WHERE id = $1 LIMIT 1;

-- name: ListArtifactsByModule :many
SELECT * FROM artifacts
WHERE module_id = $1;

-- name: ListArtifactsByModuleAndStatus :many
SELECT * FROM artifacts
WHERE module_id = $1 AND status = $2;

-- name: ListArtifactsByModuleAndType :many
SELECT * FROM artifacts
WHERE module_id = $1 AND type = $2;

-- name: UpdateArtifactStatus :exec
UPDATE artifacts
SET status = $2, approved_at = $3, rejected_at = $4
WHERE id = $1;

-- name: DeleteArtifact :exec
DELETE FROM artifacts WHERE id = $1;

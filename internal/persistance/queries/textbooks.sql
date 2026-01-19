-- name: CreateSubject :one
INSERT INTO subjects (id, name, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, url, created_at, updated_at;

-- name: CreateSubSubject :one
INSERT INTO sub_subjects (id, subject_id, name, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, subject_id, name, url, created_at, updated_at;

-- name: GetAllSubjects :many
SELECT id, name, url, created_at, updated_at
FROM subjects
ORDER BY name ASC;

-- name: GetSubjectByID :one
SELECT id, name, url, created_at, updated_at
FROM subjects
WHERE id = $1;

-- name: GetSubjectByName :one
SELECT id, name, url, created_at, updated_at
FROM subjects
WHERE name = $1;

-- name: GetSubSubjectsBySubjectID :many
SELECT id, subject_id, name, url, created_at, updated_at
FROM sub_subjects
WHERE subject_id = $1
ORDER BY name ASC;

-- name: DeleteSubject :exec
DELETE FROM subjects
WHERE id = $1;

-- name: DeleteAllSubjects :exec
DELETE FROM subjects;

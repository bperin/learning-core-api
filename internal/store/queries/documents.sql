-- name: GetDocument :one
SELECT * FROM documents
WHERE id = $1 LIMIT 1;

-- name: CreateDocument :one
INSERT INTO documents (subject_id, store_id, title, source_uri, sha256, metadata, file_name, doc_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDocumentBySubjectAndSourceURI :one
SELECT * FROM documents
WHERE subject_id = $1 AND source_uri = $2 LIMIT 1;

-- name: ListDocumentsBySubject :many
SELECT * FROM documents
WHERE subject_id = $1;

-- name: UpdateDocument :exec
UPDATE documents SET title = $2, metadata = $3, indexed_at = $4 WHERE id = $1;

-- name: DeleteDocument :exec
DELETE FROM documents WHERE id = $1;

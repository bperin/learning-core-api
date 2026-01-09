-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1 LIMIT 1;

-- name: GetDocumentsByUser :many
SELECT * FROM documents WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetDocumentsBySubject :many
SELECT * FROM documents WHERE subject_id = $1 ORDER BY created_at DESC;

-- name: GetDocumentsByRagStatus :many
SELECT * FROM documents WHERE rag_status = $1 ORDER BY created_at DESC;

-- name: ListDocuments :many
SELECT * FROM documents ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateDocument :one
INSERT INTO documents (
  filename, title, mime_type, content, storage_path, rag_status, 
  user_id, subject_id, curricular, subjects
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateDocument :one
UPDATE documents SET
  title = COALESCE($2, title),
  content = COALESCE($3, content),
  storage_path = COALESCE($4, storage_path),
  rag_status = COALESCE($5, rag_status),
  subject_id = COALESCE($6, subject_id),
  curricular = COALESCE($7, curricular),
  subjects = COALESCE($8, subjects),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateDocumentRagStatus :one
UPDATE documents SET
  rag_status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents WHERE id = $1;

-- name: SearchDocumentsByTitle :many
SELECT * FROM documents 
WHERE title ILIKE '%' || @title || '%' 
ORDER BY created_at DESC 
LIMIT @page_limit OFFSET @page_offset;

-- name: GetDocumentsBySubjects :many
SELECT * FROM documents 
WHERE subjects && $1::text[]
ORDER BY created_at DESC;

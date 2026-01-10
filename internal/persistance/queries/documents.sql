-- name: GetDocument :one
SELECT * FROM documents WHERE id = $1 LIMIT 1;

-- name: GetDocumentsByUser :many
SELECT * FROM documents WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetDocumentsByRagStatus :many
SELECT * FROM documents WHERE rag_status = $1 ORDER BY created_at DESC;

-- name: ListDocuments :many
SELECT * FROM documents ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CreateDocument :one
INSERT INTO documents (
  filename, title, mime_type, content, storage_path, storage_bucket,
  file_store_name, file_store_file_name, rag_status,
  user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateDocument :one
UPDATE documents SET
  title = COALESCE($2, title),
  content = COALESCE($3, content),
  storage_path = COALESCE($4, storage_path),
  storage_bucket = COALESCE($5, storage_bucket),
  file_store_name = COALESCE($6, file_store_name),
  file_store_file_name = COALESCE($7, file_store_file_name),
  rag_status = COALESCE($8, rag_status),
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

-- name: CreateFileSearchStore :one
INSERT INTO file_search_stores (subject_id, store_name, display_name, chunking_config)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFileSearchStore :one
SELECT * FROM file_search_stores
WHERE id = $1 LIMIT 1;

-- name: GetFileSearchStoreBySubjectID :one
SELECT * FROM file_search_stores
WHERE subject_id = $1 LIMIT 1;

-- name: UpdateFileSearchStore :one
UPDATE file_search_stores
SET display_name = $2,
    chunking_config = $3
WHERE id = $1
RETURNING *;

-- name: DeleteFileSearchStore :exec
DELETE FROM file_search_stores WHERE id = $1;

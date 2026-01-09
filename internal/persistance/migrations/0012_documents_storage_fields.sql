-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Documents: storage + file store metadata
-- =====================================================
ALTER TABLE documents
  ADD COLUMN IF NOT EXISTS storage_bucket TEXT,
  ADD COLUMN IF NOT EXISTS file_store_name TEXT,
  ADD COLUMN IF NOT EXISTS file_store_file_name TEXT;

COMMENT ON COLUMN documents.storage_bucket IS 'GCS bucket storing the source file';
COMMENT ON COLUMN documents.file_store_name IS 'File search store name for retrieval';
COMMENT ON COLUMN documents.file_store_file_name IS 'File name returned by file search store';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE documents
  DROP COLUMN IF EXISTS file_store_file_name,
  DROP COLUMN IF EXISTS file_store_name,
  DROP COLUMN IF EXISTS storage_bucket;

-- +goose StatementEnd

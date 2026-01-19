-- +goose Up
-- Add textbook_id column to documents table to link downloaded textbooks
ALTER TABLE documents
  ADD COLUMN IF NOT EXISTS textbook_id UUID REFERENCES subjects(id) ON DELETE SET NULL;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_documents_textbook ON documents(textbook_id);

-- Add comment for documentation
COMMENT ON COLUMN documents.textbook_id IS 'Reference to the textbook subject if document was downloaded from Open Textbook Library';

-- +goose Down
DROP INDEX IF EXISTS idx_documents_textbook;
ALTER TABLE documents
  DROP COLUMN IF EXISTS textbook_id;

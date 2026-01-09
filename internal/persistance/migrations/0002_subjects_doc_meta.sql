-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Update Documents with Curricular and Subjects
-- =====================================================

-- Add curricular classification to documents
ALTER TABLE documents 
ADD COLUMN IF NOT EXISTS curricular TEXT;

-- Add subjects (array of strings) to documents for broader categorization
ALTER TABLE documents 
ADD COLUMN IF NOT EXISTS subjects TEXT[] DEFAULT '{}';

-- Add comments for documentation
COMMENT ON COLUMN documents.curricular IS 'Curricular classification or framework (e.g., Common Core, IB)';
COMMENT ON COLUMN documents.subjects IS 'List of academic subjects associated with this document';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the added columns
ALTER TABLE documents DROP COLUMN IF EXISTS subjects;
ALTER TABLE documents DROP COLUMN IF EXISTS curricular;

-- +goose StatementEnd

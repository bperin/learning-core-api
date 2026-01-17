-- +goose Up
ALTER TABLE eval_items ADD COLUMN grounding_metadata JSONB;
ALTER TABLE eval_items ADD COLUMN source_document_id UUID REFERENCES documents(id);
COMMENT ON COLUMN eval_items.grounding_metadata IS 'Grounding metadata from generation (chunks, supporting text, etc.) used for eval';
COMMENT ON COLUMN eval_items.source_document_id IS 'Source document this question was generated from';

-- +goose Down
ALTER TABLE eval_items DROP COLUMN source_document_id;
ALTER TABLE eval_items DROP COLUMN grounding_metadata;

-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
  DROP CONSTRAINT IF EXISTS documents_module_id_source_uri_key;

ALTER TABLE documents
  DROP COLUMN IF EXISTS module_id,
  ADD COLUMN subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE;

ALTER TABLE documents
  ADD CONSTRAINT documents_subject_id_source_uri_key UNIQUE (subject_id, source_uri);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents
  DROP CONSTRAINT IF EXISTS documents_subject_id_source_uri_key;

ALTER TABLE documents
  DROP COLUMN IF EXISTS subject_id,
  ADD COLUMN module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE;

ALTER TABLE documents
  ADD CONSTRAINT documents_module_id_source_uri_key UNIQUE (module_id, source_uri);
-- +goose StatementEnd

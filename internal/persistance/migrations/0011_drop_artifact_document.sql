-- +goose Up
-- +goose StatementBegin

ALTER TABLE artifacts DROP COLUMN IF EXISTS document_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE artifacts
  ADD COLUMN IF NOT EXISTS document_id UUID REFERENCES documents(id) ON DELETE CASCADE;

-- +goose StatementEnd

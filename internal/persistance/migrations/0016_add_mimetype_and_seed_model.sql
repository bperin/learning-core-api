-- +goose Up
-- +goose StatementBegin

-- Add mime_type column
ALTER TABLE model_configs ADD COLUMN mime_type TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE model_configs DROP COLUMN mime_type;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- Add name column to modules table as alias for title
ALTER TABLE modules ADD COLUMN IF NOT EXISTS name TEXT;
-- Copy existing title values to name
UPDATE modules SET name = title WHERE name IS NULL OR name = '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove name column from modules table
ALTER TABLE modules DROP COLUMN IF EXISTS name;

-- +goose StatementEnd
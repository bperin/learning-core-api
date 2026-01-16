-- +goose Up
-- +goose StatementBegin

ALTER TYPE generation_type ADD VALUE IF NOT EXISTS 'QUESTIONS';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Enum values cannot be removed safely; no-op rollback.
SELECT 1;

-- +goose StatementEnd

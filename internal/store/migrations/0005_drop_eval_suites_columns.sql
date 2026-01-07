-- +goose Up
-- +goose StatementBegin
ALTER TABLE eval_suites DROP COLUMN thresholds;
ALTER TABLE eval_suites DROP COLUMN eval_types;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE eval_suites ADD COLUMN thresholds JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE eval_suites ADD COLUMN eval_types eval_type[] NOT NULL DEFAULT '{}';
-- +goose StatementEnd
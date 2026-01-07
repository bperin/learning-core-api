-- +goose Up
-- +goose StatementBegin
ALTER TABLE eval_runs
  ALTER COLUMN suite_id SET NOT NULL;

ALTER TABLE eval_results
  ALTER COLUMN rule_id SET NOT NULL;

ALTER TABLE eval_results
  DROP COLUMN IF EXISTS type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE eval_results
  ADD COLUMN type eval_type NOT NULL DEFAULT 'SCHEMA_VALIDATION';

ALTER TABLE eval_results
  ALTER COLUMN type DROP DEFAULT;

ALTER TABLE eval_results
  ALTER COLUMN rule_id DROP NOT NULL;

ALTER TABLE eval_runs
  ALTER COLUMN suite_id DROP NOT NULL;
-- +goose StatementEnd

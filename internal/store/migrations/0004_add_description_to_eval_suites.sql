-- +goose Up
-- +goose StatementBegin
ALTER TABLE eval_suites
  ADD COLUMN description TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE eval_suites
  DROP COLUMN IF EXISTS description;
-- +goose StatementEnd

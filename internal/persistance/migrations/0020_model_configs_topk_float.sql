-- +goose Up
-- +goose StatementBegin

ALTER TABLE model_configs
  ALTER COLUMN top_k TYPE REAL
  USING top_k::REAL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE model_configs
  ALTER COLUMN top_k TYPE INTEGER
  USING ROUND(top_k)::INTEGER;

-- +goose StatementEnd

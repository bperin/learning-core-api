-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subjects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, name)
);

DROP INDEX IF EXISTS ux_file_search_store_module;
ALTER TABLE file_search_stores
  DROP COLUMN IF EXISTS module_id,
  ADD COLUMN subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX IF NOT EXISTS ux_file_search_store_subject
  ON file_search_stores(subject_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_file_search_store_subject;
ALTER TABLE file_search_stores
  DROP COLUMN IF EXISTS subject_id,
  ADD COLUMN module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX IF NOT EXISTS ux_file_search_store_module
  ON file_search_stores(module_id);

DROP TABLE IF EXISTS subjects;
-- +goose StatementEnd

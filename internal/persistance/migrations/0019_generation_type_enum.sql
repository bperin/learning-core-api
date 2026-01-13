-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Generation Types
-- =====================================================
CREATE TYPE generation_type AS ENUM ('CLASSIFICATION');

-- =====================================================
-- Prompt Templates
-- =====================================================
ALTER TABLE prompt_templates RENAME COLUMN key TO generation_type;
UPDATE prompt_templates SET generation_type = 'CLASSIFICATION';
ALTER TABLE prompt_templates
  ALTER COLUMN generation_type TYPE generation_type USING generation_type::generation_type;
ALTER TABLE prompt_templates
  ALTER COLUMN generation_type SET NOT NULL;

ALTER TABLE prompt_templates DROP CONSTRAINT IF EXISTS prompt_templates_key_version_key;
ALTER TABLE prompt_templates DROP CONSTRAINT IF EXISTS prompt_templates_generation_type_version_key;
ALTER TABLE prompt_templates
  ADD CONSTRAINT prompt_templates_generation_type_version_key UNIQUE (generation_type, version);

DROP INDEX IF EXISTS idx_prompt_templates_key_active;
CREATE INDEX IF NOT EXISTS idx_prompt_templates_generation_type_active
  ON prompt_templates(generation_type, is_active);

COMMENT ON COLUMN prompt_templates.generation_type IS 'Generation type this prompt supports';

-- =====================================================
-- Schema Templates
-- =====================================================
ALTER TABLE schema_templates RENAME COLUMN schema_type TO generation_type;
UPDATE schema_templates SET generation_type = 'CLASSIFICATION';
ALTER TABLE schema_templates
  ALTER COLUMN generation_type TYPE generation_type USING generation_type::generation_type;
ALTER TABLE schema_templates
  ALTER COLUMN generation_type SET NOT NULL;

ALTER TABLE schema_templates DROP CONSTRAINT IF EXISTS schema_templates_schema_type_version_key;
ALTER TABLE schema_templates DROP CONSTRAINT IF EXISTS schema_templates_generation_type_version_key;
ALTER TABLE schema_templates
  ADD CONSTRAINT schema_templates_generation_type_version_key UNIQUE (generation_type, version);

DROP INDEX IF EXISTS idx_schema_templates_type;
CREATE INDEX IF NOT EXISTS idx_schema_templates_generation_type
  ON schema_templates(generation_type);

COMMENT ON COLUMN schema_templates.generation_type IS 'Generation type this schema supports';

-- =====================================================
-- Artifacts
-- =====================================================
ALTER TABLE artifacts ADD COLUMN IF NOT EXISTS generation_type generation_type;
COMMENT ON COLUMN artifacts.generation_type IS 'Generation type associated with the prompt/schema output';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE artifacts DROP COLUMN IF EXISTS generation_type;

DROP INDEX IF EXISTS idx_schema_templates_generation_type;
ALTER TABLE schema_templates DROP CONSTRAINT IF EXISTS schema_templates_generation_type_version_key;
ALTER TABLE schema_templates
  ALTER COLUMN generation_type TYPE TEXT USING generation_type::text;
ALTER TABLE schema_templates RENAME COLUMN generation_type TO schema_type;
ALTER TABLE schema_templates
  ADD CONSTRAINT schema_templates_schema_type_version_key UNIQUE (schema_type, version);
CREATE INDEX IF NOT EXISTS idx_schema_templates_type ON schema_templates(schema_type);

DROP INDEX IF EXISTS idx_prompt_templates_generation_type_active;
ALTER TABLE prompt_templates DROP CONSTRAINT IF EXISTS prompt_templates_generation_type_version_key;
ALTER TABLE prompt_templates
  ALTER COLUMN generation_type TYPE TEXT USING generation_type::text;
ALTER TABLE prompt_templates RENAME COLUMN generation_type TO key;
ALTER TABLE prompt_templates
  ADD CONSTRAINT prompt_templates_key_version_key UNIQUE (key, version);
CREATE INDEX IF NOT EXISTS idx_prompt_templates_key_active ON prompt_templates(key, is_active);

DROP TYPE IF EXISTS generation_type;

-- +goose StatementEnd

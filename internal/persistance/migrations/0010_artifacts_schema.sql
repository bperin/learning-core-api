-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Artifacts: immutability + provenance fields
-- =====================================================

-- Remove updated_at tracking (artifacts are append-only)
DROP TRIGGER IF EXISTS update_artifacts_updated_at ON artifacts;
DROP INDEX IF EXISTS idx_artifacts_updated_at;
ALTER TABLE artifacts DROP COLUMN IF EXISTS updated_at;

-- Rename ambiguous columns
ALTER TABLE artifacts RENAME COLUMN json TO output_json;
ALTER TABLE artifacts RENAME COLUMN user_id TO reviewer_id;

-- Add reproducibility fields
ALTER TABLE artifacts
  ADD COLUMN IF NOT EXISTS prompt_template_id UUID REFERENCES prompt_templates(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS schema_template_id UUID REFERENCES schema_templates(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS model_params JSONB,
  ADD COLUMN IF NOT EXISTS prompt_render TEXT;

COMMENT ON COLUMN artifacts.output_json IS 'Structured output payload from the model';
COMMENT ON COLUMN artifacts.reviewer_id IS 'User who reviewed or created the artifact';
COMMENT ON COLUMN artifacts.prompt_template_id IS 'Prompt template used for generation';
COMMENT ON COLUMN artifacts.schema_template_id IS 'Schema template used to validate output';
COMMENT ON COLUMN artifacts.model_params IS 'Model parameters (temperature, top_p, etc.)';
COMMENT ON COLUMN artifacts.prompt_render IS 'Rendered prompt text sent to the model';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE artifacts
  DROP COLUMN IF EXISTS prompt_render,
  DROP COLUMN IF EXISTS model_params,
  DROP COLUMN IF EXISTS schema_template_id,
  DROP COLUMN IF EXISTS prompt_template_id;

ALTER TABLE artifacts RENAME COLUMN reviewer_id TO user_id;
ALTER TABLE artifacts RENAME COLUMN output_json TO json;

ALTER TABLE artifacts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();
CREATE INDEX IF NOT EXISTS idx_artifacts_updated_at ON artifacts(updated_at);

DROP TRIGGER IF EXISTS update_artifacts_updated_at ON artifacts;
CREATE TRIGGER update_artifacts_updated_at
  BEFORE UPDATE ON artifacts
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

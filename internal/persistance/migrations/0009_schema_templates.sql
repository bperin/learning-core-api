-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Schema Templates
-- =====================================================
CREATE TABLE schema_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  schema_type TEXT NOT NULL,
  version INTEGER NOT NULL,
  schema_json JSONB NOT NULL,
  subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
  curriculum_id UUID REFERENCES curricula(id) ON DELETE SET NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  locked_at TIMESTAMPTZ,
  UNIQUE (schema_type, version)
);

CREATE INDEX idx_schema_templates_type ON schema_templates(schema_type);
CREATE INDEX idx_schema_templates_active ON schema_templates(is_active);
CREATE INDEX idx_schema_templates_subject ON schema_templates(subject_id);
CREATE INDEX idx_schema_templates_curriculum ON schema_templates(curriculum_id);

COMMENT ON COLUMN schema_templates.schema_type IS 'Schema purpose (eval_generation, intent_extraction, etc.)';
COMMENT ON COLUMN schema_templates.schema_json IS 'JSON schema defining expected AI output';
COMMENT ON COLUMN schema_templates.locked_at IS 'Timestamp when version becomes immutable';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS schema_templates;

-- +goose StatementEnd

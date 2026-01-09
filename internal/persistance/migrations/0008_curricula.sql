-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Curricula
-- =====================================================
CREATE TABLE curricula (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
  parent_id UUID REFERENCES curricula(id) ON DELETE SET NULL,
  label TEXT NOT NULL,
  code TEXT,
  description TEXT,
  order_index INTEGER NOT NULL DEFAULT 0,
  grade_level TEXT,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_curricula_subject ON curricula(subject_id);
CREATE INDEX idx_curricula_parent ON curricula(parent_id);

DROP TRIGGER IF EXISTS update_curricula_updated_at ON curricula;
CREATE TRIGGER update_curricula_updated_at
  BEFORE UPDATE ON curricula
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

COMMENT ON COLUMN curricula.subject_id IS 'Subject this curriculum belongs to';
COMMENT ON COLUMN curricula.parent_id IS 'Parent curriculum node; NULL is root';
COMMENT ON COLUMN curricula.label IS 'Human-readable curriculum label';
COMMENT ON COLUMN curricula.code IS 'Optional standardized curriculum code';
COMMENT ON COLUMN curricula.description IS 'Optional curriculum description';
COMMENT ON COLUMN curricula.order_index IS 'Sort order within parent curriculum';
COMMENT ON COLUMN curricula.grade_level IS 'Grade or level alignment';
COMMENT ON COLUMN curricula.is_active IS 'Whether curriculum unit is active';

-- =====================================================
-- Documents: add curriculum relationship
-- =====================================================
ALTER TABLE documents ADD COLUMN IF NOT EXISTS curriculum_id UUID REFERENCES curricula(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_documents_curriculum ON documents(curriculum_id);

COMMENT ON COLUMN documents.curriculum_id IS 'Curriculum unit associated with this document';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE documents DROP COLUMN IF EXISTS curriculum_id;
DROP INDEX IF EXISTS idx_documents_curriculum;

DROP TRIGGER IF EXISTS update_curricula_updated_at ON curricula;
DROP INDEX IF EXISTS idx_curricula_parent;
DROP INDEX IF EXISTS idx_curricula_subject;
DROP TABLE IF EXISTS curricula;

-- +goose StatementEnd

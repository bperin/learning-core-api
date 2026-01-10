-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Taxonomy Nodes
-- =====================================================
CREATE TABLE IF NOT EXISTS taxonomy_nodes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  description TEXT,
  parent_id UUID REFERENCES taxonomy_nodes(id) ON DELETE CASCADE,
  path TEXT NOT NULL,
  depth INTEGER NOT NULL,
  state TEXT NOT NULL CHECK (state IN ('ai_generated', 'approved', 'rejected')),
  confidence FLOAT CHECK (confidence >= 0 AND confidence <= 1),
  source_document_id UUID REFERENCES documents(id),
  version INTEGER NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_by UUID REFERENCES users(id),
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (path, version)
);

CREATE INDEX IF NOT EXISTS idx_taxonomy_nodes_path ON taxonomy_nodes(path);
CREATE INDEX IF NOT EXISTS idx_taxonomy_nodes_state ON taxonomy_nodes(state);
CREATE INDEX IF NOT EXISTS idx_taxonomy_nodes_active ON taxonomy_nodes(is_active);

-- =====================================================
-- Document Taxonomy Links
-- =====================================================
CREATE TABLE IF NOT EXISTS document_taxonomy_links (
  document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
  taxonomy_node_id UUID REFERENCES taxonomy_nodes(id) ON DELETE CASCADE,
  confidence FLOAT CHECK (confidence >= 0 AND confidence <= 1),
  state TEXT NOT NULL CHECK (state IN ('ai_generated', 'approved', 'rejected')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMPTZ,
  PRIMARY KEY (document_id, taxonomy_node_id)
);

CREATE INDEX IF NOT EXISTS idx_document_taxonomy_links_state ON document_taxonomy_links(state);

-- =====================================================
-- Drop legacy subject/curriculum fields + rubric
-- =====================================================
ALTER TABLE evals DROP COLUMN IF EXISTS rubric;
ALTER TABLE evals DROP COLUMN IF EXISTS subject_id;
DROP INDEX IF EXISTS idx_evals_subject;

ALTER TABLE documents DROP COLUMN IF EXISTS subject_id;
ALTER TABLE documents DROP COLUMN IF EXISTS curriculum_id;
ALTER TABLE documents DROP COLUMN IF EXISTS curricular;
ALTER TABLE documents DROP COLUMN IF EXISTS subjects;
DROP INDEX IF EXISTS idx_documents_subject;
DROP INDEX IF EXISTS idx_documents_curriculum;

ALTER TABLE schema_templates DROP COLUMN IF EXISTS subject_id;
ALTER TABLE schema_templates DROP COLUMN IF EXISTS curriculum_id;
DROP INDEX IF EXISTS idx_schema_templates_subject;
DROP INDEX IF EXISTS idx_schema_templates_curriculum;

DROP TABLE IF EXISTS curricula;
DROP TABLE IF EXISTS subjects;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- =====================================================
-- Restore legacy subjects/curricula + rubric
-- =====================================================
CREATE TABLE IF NOT EXISTS subjects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  description TEXT,
  user_id UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subjects_user ON subjects(user_id);

CREATE TABLE IF NOT EXISTS curricula (
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

CREATE INDEX IF NOT EXISTS idx_curricula_subject ON curricula(subject_id);

ALTER TABLE documents ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS curriculum_id UUID REFERENCES curricula(id) ON DELETE SET NULL;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS curricular TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS subjects TEXT[] DEFAULT '{}';

CREATE INDEX IF NOT EXISTS idx_documents_subject ON documents(subject_id);
CREATE INDEX IF NOT EXISTS idx_documents_curriculum ON documents(curriculum_id);

ALTER TABLE evals ADD COLUMN IF NOT EXISTS rubric JSONB;
ALTER TABLE evals ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_evals_subject ON evals(subject_id);

ALTER TABLE schema_templates ADD COLUMN IF NOT EXISTS subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL;
ALTER TABLE schema_templates ADD COLUMN IF NOT EXISTS curriculum_id UUID REFERENCES curricula(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_schema_templates_subject ON schema_templates(subject_id);
CREATE INDEX IF NOT EXISTS idx_schema_templates_curriculum ON schema_templates(curriculum_id);

DROP TABLE IF EXISTS document_taxonomy_links;
DROP TABLE IF EXISTS taxonomy_nodes;

-- +goose StatementEnd

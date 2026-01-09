-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Reset legacy schema (Postgres)
-- =====================================================
DROP TABLE IF EXISTS eval_item_reviews CASCADE;
DROP TABLE IF EXISTS artifacts CASCADE;
DROP TABLE IF EXISTS user_answers CASCADE;
DROP TABLE IF EXISTS test_attempts CASCADE;
DROP TABLE IF EXISTS prompt_templates CASCADE;
DROP TABLE IF EXISTS eval_items CASCADE;
DROP TABLE IF EXISTS evals CASCADE;
DROP TABLE IF EXISTS documents CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- =====================================================
-- Users / Subjects
-- =====================================================
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  is_admin BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE subjects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  description TEXT,
  user_id UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_subjects_user ON subjects(user_id);

-- =====================================================
-- Documents
-- =====================================================
CREATE TABLE documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  filename TEXT NOT NULL,
  title TEXT,
  mime_type TEXT,
  content TEXT,
  storage_path TEXT,
  rag_status TEXT NOT NULL DEFAULT 'PENDING',
  user_id UUID NOT NULL REFERENCES users(id),
  subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_documents_user ON documents(user_id);
CREATE INDEX idx_documents_subject ON documents(subject_id);

-- =====================================================
-- Evals / Items
-- =====================================================
CREATE TABLE evals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL DEFAULT 'draft',
  difficulty TEXT,
  instructions TEXT,
  rubric JSONB,
  subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
  user_id UUID NOT NULL REFERENCES users(id),
  published_at TIMESTAMPTZ,
  archived_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_evals_user ON evals(user_id);
CREATE INDEX idx_evals_subject ON evals(subject_id);

CREATE TABLE eval_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  eval_id UUID NOT NULL REFERENCES evals(id) ON DELETE CASCADE,
  prompt TEXT NOT NULL,
  options TEXT[] NOT NULL,
  correct_idx INTEGER NOT NULL,
  hint TEXT,
  explanation TEXT,
  metadata JSONB
);

CREATE INDEX idx_eval_items_eval ON eval_items(eval_id);

-- =====================================================
-- Prompt Templates
-- =====================================================
CREATE TABLE prompt_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  is_active BOOLEAN NOT NULL DEFAULT false,
  title TEXT NOT NULL,
  description TEXT,
  template TEXT NOT NULL,
  metadata JSONB,
  created_by TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (key, version)
);

CREATE INDEX idx_prompt_templates_key_active ON prompt_templates(key, is_active);

-- =====================================================
-- Attempts / Answers
-- =====================================================
CREATE TABLE test_attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  eval_id UUID NOT NULL REFERENCES evals(id),
  score INTEGER NOT NULL DEFAULT 0,
  total INTEGER NOT NULL,
  percentage REAL,
  total_time INTEGER,
  feedback JSONB,
  summary TEXT,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE INDEX idx_test_attempts_user ON test_attempts(user_id);
CREATE INDEX idx_test_attempts_eval ON test_attempts(eval_id);

CREATE TABLE user_answers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  attempt_id UUID NOT NULL REFERENCES test_attempts(id) ON DELETE CASCADE,
  eval_item_id UUID NOT NULL REFERENCES eval_items(id),
  selected_idx INTEGER NOT NULL,
  is_correct BOOLEAN NOT NULL,
  time_spent INTEGER,
  hints_used INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_answers_attempt ON user_answers(attempt_id);
CREATE INDEX idx_user_answers_eval_item ON user_answers(eval_item_id);

-- =====================================================
-- Artifacts
-- =====================================================
CREATE TABLE artifacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'READY',
  document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
  eval_id UUID REFERENCES evals(id) ON DELETE CASCADE,
  eval_item_id UUID REFERENCES eval_items(id) ON DELETE CASCADE,
  attempt_id UUID REFERENCES test_attempts(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  text TEXT,
  json JSONB,
  model TEXT,
  prompt TEXT,
  input_hash TEXT,
  meta JSONB,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_artifacts_type ON artifacts(type);

-- =====================================================
-- Eval Item Reviews
-- =====================================================
CREATE TABLE eval_item_reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  eval_item_id UUID NOT NULL REFERENCES eval_items(id) ON DELETE CASCADE,
  reviewer_id UUID NOT NULL REFERENCES users(id),
  verdict TEXT NOT NULL,
  reasons TEXT[] NOT NULL,
  comments TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_eval_item_reviews_item ON eval_item_reviews(eval_item_id);
CREATE INDEX idx_eval_item_reviews_reviewer ON eval_item_reviews(reviewer_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS eval_item_reviews;
DROP TABLE IF EXISTS artifacts;
DROP TABLE IF EXISTS user_answers;
DROP TABLE IF EXISTS test_attempts;
DROP TABLE IF EXISTS prompt_templates;
DROP TABLE IF EXISTS eval_items;
DROP TABLE IF EXISTS evals;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

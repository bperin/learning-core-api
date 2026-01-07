-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Extensions
-- =====================================================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- Enums
-- =====================================================
DO $$ BEGIN
  CREATE TYPE user_role_type AS ENUM ('ADMIN', 'INSTRUCTOR', 'LEARNER');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE artifact_status AS ENUM (
    'PENDING_GENERATION',
    'PENDING_EVAL',
    'APPROVED',
    'REJECTED',
    'RETIRED'
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE artifact_type AS ENUM (
    'MCQ_ITEM',
    'FLASHCARD',
    'STUDY_PLAN'
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE run_status AS ENUM (
    'PENDING',
    'RUNNING',
    'SUCCEEDED',
    'FAILED'
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE eval_type AS ENUM (
    'SCHEMA_VALIDATION',
    'GROUNDEDNESS',
    'ANSWER_CORRECTNESS',
    'DISTRACTOR_QUALITY',
    'DIFFICULTY_CALIBRATION',
    'CONCEPT_ALIGNMENT'
  );
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- =====================================================
-- Tenants / Users / Roles
-- =====================================================
CREATE TABLE IF NOT EXISTS tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  email TEXT NOT NULL,
  display_name TEXT,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (tenant_id, email)
);

CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);

CREATE TABLE IF NOT EXISTS user_roles (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role user_role_type NOT NULL,
  granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role)
);

-- =====================================================
-- Modules
-- =====================================================
CREATE TABLE IF NOT EXISTS modules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_modules_tenant ON modules(tenant_id);

-- =====================================================
-- File Search Stores (1 per module)
-- =====================================================
CREATE TABLE IF NOT EXISTS file_search_stores (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  store_name TEXT NOT NULL UNIQUE,
  display_name TEXT,
  chunking_config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_file_search_store_module
  ON file_search_stores(module_id);

-- =====================================================
-- Documents
-- =====================================================
CREATE TABLE IF NOT EXISTS documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  store_id UUID NOT NULL REFERENCES file_search_stores(id) ON DELETE CASCADE,
  title TEXT,
  source_uri TEXT NOT NULL,
  sha256 TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  file_name TEXT,
  doc_name TEXT,
  indexed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (module_id, source_uri)
);

CREATE INDEX IF NOT EXISTS idx_documents_module ON documents(module_id);

-- =====================================================
-- Prompt Versions
-- =====================================================
CREATE TABLE IF NOT EXISTS prompt_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  version TEXT NOT NULL,
  content TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (name, version)
);

-- =====================================================
-- Generation Runs
-- =====================================================
CREATE TABLE IF NOT EXISTS generation_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  agent_name TEXT NOT NULL,
  agent_version TEXT NOT NULL,
  model TEXT NOT NULL,
  model_params JSONB NOT NULL DEFAULT '{}'::jsonb,
  prompt_id UUID REFERENCES prompt_versions(id),
  store_name TEXT NOT NULL,
  metadata_filter JSONB NOT NULL DEFAULT '{}'::jsonb,
  status run_status NOT NULL DEFAULT 'PENDING',
  input_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  output_payload JSONB,
  error JSONB,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_generation_runs_module
  ON generation_runs(module_id, created_at DESC);

-- =====================================================
-- Artifacts (Eval-gated learning content)
-- =====================================================
CREATE TABLE IF NOT EXISTS artifacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  generation_run_id UUID NOT NULL REFERENCES generation_runs(id) ON DELETE CASCADE,
  type artifact_type NOT NULL,
  status artifact_status NOT NULL DEFAULT 'PENDING_EVAL',
  schema_version TEXT NOT NULL,
  difficulty TEXT,
  tags TEXT[] NOT NULL DEFAULT '{}',
  artifact_payload JSONB NOT NULL,
  grounding JSONB NOT NULL DEFAULT '{}'::jsonb,
  evidence_version TEXT,
  approved_at TIMESTAMPTZ,
  rejected_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_artifacts_pending_eval
  ON artifacts(status, created_at DESC)
  WHERE status = 'PENDING_EVAL';

-- =====================================================
-- Eval Suites / Runs / Results
-- =====================================================
CREATE TABLE IF NOT EXISTS eval_suites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  eval_types eval_type[] NOT NULL,
  thresholds JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS eval_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
  generation_run_id UUID REFERENCES generation_runs(id),
  suite_id UUID REFERENCES eval_suites(id),
  judge_model TEXT NOT NULL,
  judge_params JSONB NOT NULL DEFAULT '{}'::jsonb,
  status run_status NOT NULL DEFAULT 'PENDING',
  overall_pass BOOLEAN,
  overall_score REAL,
  error JSONB,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS eval_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  eval_run_id UUID NOT NULL REFERENCES eval_runs(id) ON DELETE CASCADE,
  type eval_type NOT NULL,
  pass BOOLEAN NOT NULL,
  score REAL,
  details JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- Sessions / Attempts (Learner runtime)
-- =====================================================
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  mastery_state JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_user
  ON sessions(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
  artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE RESTRICT,
  is_correct BOOLEAN NOT NULL,
  user_answer JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_attempts_session
  ON attempts(session_id, created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attempts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS eval_results;
DROP TABLE IF EXISTS eval_runs;
DROP TABLE IF EXISTS eval_suites;
DROP TABLE IF EXISTS artifacts;
DROP TABLE IF EXISTS generation_runs;
DROP TABLE IF EXISTS prompt_versions;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS file_search_stores;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
DROP TYPE IF EXISTS eval_type;
DROP TYPE IF EXISTS run_status;
DROP TYPE IF EXISTS artifact_type;
DROP TYPE IF EXISTS artifact_status;
DROP TYPE IF EXISTS user_role_type;
-- +goose StatementEnd

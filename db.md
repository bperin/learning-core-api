ALTER TABLE attempts
ADD COLUMN IF NOT EXISTS tenant_id UUID
  REFERENCES tenants(id) ON DELETE CASCADE;

UPDATE attempts a
SET tenant_id = s.tenant_id
FROM sessions s
WHERE a.session_id = s.id
  AND a.tenant_id IS NULL;

ALTER TABLE attempts
ALTER COLUMN tenant_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_attempts_tenant
  ON attempts (tenant_id);
2️⃣ Remaining tables (fully written, final form)
Below are final versions of tables that were previously partial or implicit.

Documents (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS documents (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id     UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
  store_id      UUID NOT NULL REFERENCES file_search_stores(id) ON DELETE CASCADE,

  title         TEXT,
  source_uri    TEXT NOT NULL,
  sha256        TEXT,

  metadata      JSONB NOT NULL DEFAULT '{}'::jsonb,

  file_name     TEXT,
  doc_name      TEXT,
  indexed_at    TIMESTAMPTZ,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (module_id, source_uri)
);

CREATE INDEX IF NOT EXISTS idx_documents_module
  ON documents (module_id);
Generation runs (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS generation_runs (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id      UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,

  agent_name     TEXT NOT NULL,
  agent_version  TEXT NOT NULL,

  model          TEXT NOT NULL,
  model_params   JSONB NOT NULL DEFAULT '{}'::jsonb,

  prompt_id      UUID REFERENCES prompt_versions(id),

  store_name     TEXT NOT NULL,
  metadata_filter JSONB NOT NULL DEFAULT '{}'::jsonb,

  status         run_status NOT NULL DEFAULT 'PENDING',
  input_payload  JSONB NOT NULL DEFAULT '{}'::jsonb,
  output_payload JSONB,
  error          JSONB,

  started_at     TIMESTAMPTZ,
  finished_at    TIMESTAMPTZ,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_generation_runs_module
  ON generation_runs (module_id, created_at DESC);
Artifacts (final, eval-gated learning content)
sql
Copy code
CREATE TABLE IF NOT EXISTS artifacts (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  module_id       UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,

  generation_run_id UUID NOT NULL REFERENCES generation_runs(id) ON DELETE CASCADE,

  type            artifact_type NOT NULL,
  status          artifact_status NOT NULL DEFAULT 'PENDING_EVAL',

  schema_version  TEXT NOT NULL,
  difficulty      TEXT,
  tags            TEXT[] NOT NULL DEFAULT '{}',

  artifact_payload JSONB NOT NULL,
  grounding        JSONB NOT NULL DEFAULT '{}'::jsonb,

  evidence_version TEXT,          -- hash or store revision
  approved_at     TIMESTAMPTZ,
  rejected_at     TIMESTAMPTZ,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_artifacts_module_status
  ON artifacts (module_id, status);

CREATE INDEX IF NOT EXISTS idx_artifacts_type
  ON artifacts (type);
Eval suites (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS eval_suites (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name          TEXT NOT NULL UNIQUE,
  description   TEXT NOT NULL DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
Eval rules (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS eval_rules (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  suite_id      UUID NOT NULL REFERENCES eval_suites(id) ON DELETE CASCADE,

  eval_type     eval_type NOT NULL,
  min_score     REAL,
  max_score     REAL,
  weight        REAL NOT NULL DEFAULT 1.0,
  hard_fail     BOOLEAN NOT NULL DEFAULT FALSE,
  params        JSONB NOT NULL DEFAULT '{}'::jsonb,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (suite_id, eval_type)
);
Eval runs (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS eval_runs (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  artifact_id     UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
  generation_run_id UUID REFERENCES generation_runs(id),

  suite_id        UUID NOT NULL REFERENCES eval_suites(id),
  judge_model     TEXT NOT NULL,
  judge_params    JSONB NOT NULL DEFAULT '{}'::jsonb,

  status          run_status NOT NULL DEFAULT 'PENDING',
  overall_pass    BOOLEAN,
  overall_score   REAL,
  error           JSONB,

  started_at      TIMESTAMPTZ,
  finished_at     TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_eval_runs_artifact
  ON eval_runs (artifact_id, created_at DESC);
Eval results (final)
sql
Copy code
CREATE TABLE IF NOT EXISTS eval_results (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  eval_run_id   UUID NOT NULL REFERENCES eval_runs(id) ON DELETE CASCADE,
  rule_id       UUID NOT NULL REFERENCES eval_rules(id) ON DELETE CASCADE,

  pass          BOOLEAN NOT NULL,
  score         REAL,
  details       JSONB NOT NULL DEFAULT '{}'::jsonb,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (eval_run_id, rule_id)
);

CREATE INDEX IF NOT EXISTS idx_eval_results_run
  ON eval_results (eval_run_id);
3️⃣ Final invariants (this is the important part)
Ownership rules (never violate these)
Tenants own modules

Modules own documents

Documents define evidence

Evidence → generation → artifacts

Artifacts → evals → approved bank

Users interact only with approved artifacts

Why shared DB is still correct
Tenancy stops at module/session boundary

Learning correctness stays tenant-agnostic

Evals are globally meaningful

Jobs never need auth logic

This is exactly how you scale later without rewrites.

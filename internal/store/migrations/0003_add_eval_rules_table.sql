-- +goose Up
-- +goose StatementBegin

-- -------------------------------------------------
-- Eval Rules (normalized thresholds / policy)
-- -------------------------------------------------
CREATE TABLE IF NOT EXISTS eval_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  suite_id UUID NOT NULL REFERENCES eval_suites(id) ON DELETE CASCADE,

  eval_type eval_type NOT NULL,

  min_score REAL,
  max_score REAL,
  weight REAL NOT NULL DEFAULT 1.0,
  hard_fail BOOLEAN NOT NULL DEFAULT FALSE,

  params JSONB NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (suite_id, eval_type)
);

CREATE INDEX IF NOT EXISTS idx_eval_rules_suite
  ON eval_rules(suite_id);

-- -------------------------------------------------
-- Update eval_results to reference rules
-- -------------------------------------------------
ALTER TABLE eval_results
  ADD COLUMN IF NOT EXISTS rule_id UUID;

ALTER TABLE eval_results
  ADD CONSTRAINT eval_results_rule_id_fkey
  FOREIGN KEY (rule_id) REFERENCES eval_rules(id) ON DELETE CASCADE;

-- NOTE:
-- Existing rows (if any) will have rule_id NULL.
-- New evals must populate rule_id.
-- You can backfill later if needed.

CREATE INDEX IF NOT EXISTS idx_eval_results_rule
  ON eval_results(rule_id);

-- -------------------------------------------------
-- (Optional but recommended) enforce 1 result per rule per run
-- -------------------------------------------------
DO $$ BEGIN
  ALTER TABLE eval_results
    ADD CONSTRAINT eval_results_unique_rule_per_run
    UNIQUE (eval_run_id, rule_id);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE eval_results
  DROP CONSTRAINT IF EXISTS eval_results_unique_rule_per_run;

ALTER TABLE eval_results
  DROP CONSTRAINT IF EXISTS eval_results_rule_id_fkey;

ALTER TABLE eval_results
  DROP COLUMN IF EXISTS rule_id;

DROP TABLE IF EXISTS eval_rules;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Model Configurations
-- =====================================================
CREATE TABLE model_configs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  provider TEXT NOT NULL,
  model_name TEXT NOT NULL,
  temperature REAL,
  max_tokens INTEGER,
  top_p REAL,
  top_k INTEGER,
  is_active BOOLEAN NOT NULL DEFAULT false,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (key, version)
);

CREATE INDEX idx_model_configs_key_active ON model_configs(key, is_active);

-- =====================================================
-- System Instructions
-- =====================================================
CREATE TABLE system_instructions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  text TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT false,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (key, version)
);

CREATE INDEX idx_system_instructions_key_active ON system_instructions(key, is_active);

-- =====================================================
-- Chunking Configurations
-- =====================================================
CREATE TABLE chunking_configs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  chunk_size INTEGER NOT NULL,
  chunk_overlap INTEGER NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT false,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (key, version)
);

CREATE INDEX idx_chunking_configs_key_active ON chunking_configs(key, is_active);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chunking_configs;
DROP TABLE IF EXISTS system_instructions;
DROP TABLE IF EXISTS model_configs;
-- +goose StatementEnd

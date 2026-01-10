-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Model Configurations (remove key/provider)
-- =====================================================
ALTER TABLE model_configs DROP CONSTRAINT IF EXISTS model_configs_key_version_key;
DROP INDEX IF EXISTS idx_model_configs_key_active;
ALTER TABLE model_configs DROP COLUMN IF EXISTS key;
ALTER TABLE model_configs DROP COLUMN IF EXISTS provider;

-- =====================================================
-- System Instructions (remove key)
-- =====================================================
ALTER TABLE system_instructions DROP CONSTRAINT IF EXISTS system_instructions_key_version_key;
DROP INDEX IF EXISTS idx_system_instructions_key_active;
ALTER TABLE system_instructions DROP COLUMN IF EXISTS key;

-- =====================================================
-- Chunking Configurations (remove key)
-- =====================================================
ALTER TABLE chunking_configs DROP CONSTRAINT IF EXISTS chunking_configs_key_version_key;
DROP INDEX IF EXISTS idx_chunking_configs_key_active;
ALTER TABLE chunking_configs DROP COLUMN IF EXISTS key;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- =====================================================
-- Model Configurations (restore key/provider)
-- =====================================================
ALTER TABLE model_configs ADD COLUMN key TEXT NOT NULL DEFAULT 'default';
ALTER TABLE model_configs ADD COLUMN provider TEXT NOT NULL DEFAULT 'gemini';
CREATE UNIQUE INDEX IF NOT EXISTS idx_model_configs_key_active ON model_configs(key, is_active);
ALTER TABLE model_configs ADD CONSTRAINT model_configs_key_version_key UNIQUE (key, version);

-- =====================================================
-- System Instructions (restore key)
-- =====================================================
ALTER TABLE system_instructions ADD COLUMN key TEXT NOT NULL DEFAULT 'default';
CREATE UNIQUE INDEX IF NOT EXISTS idx_system_instructions_key_active ON system_instructions(key, is_active);
ALTER TABLE system_instructions ADD CONSTRAINT system_instructions_key_version_key UNIQUE (key, version);

-- =====================================================
-- Chunking Configurations (restore key)
-- =====================================================
ALTER TABLE chunking_configs ADD COLUMN key TEXT NOT NULL DEFAULT 'default';
CREATE UNIQUE INDEX IF NOT EXISTS idx_chunking_configs_key_active ON chunking_configs(key, is_active);
ALTER TABLE chunking_configs ADD CONSTRAINT chunking_configs_key_version_key UNIQUE (key, version);

-- +goose StatementEnd

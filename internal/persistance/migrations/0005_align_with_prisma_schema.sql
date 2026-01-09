-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Create Enums
-- =====================================================

-- Create ArtifactType enum
CREATE TYPE artifact_type AS ENUM (
  'INTENTS',
  'PLAN', 
  'EVAL',
  'EVAL_ITEM',
  'PROMPT',
  'PROMPT_POLICY',
  'QUALITY_METRICS',
  'HINT',
  'SUMMARY',
  'OUTLINE',
  'OTHER'
);

-- Create ReviewVerdict enum  
CREATE TYPE review_verdict AS ENUM (
  'APPROVED',
  'REJECTED', 
  'NEEDS_REVISION'
);

-- =====================================================
-- Update Users Table
-- =====================================================

-- Add missing boolean role fields if they don't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_learner BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_teacher BOOLEAN NOT NULL DEFAULT false;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_is_learner ON users(is_learner);
CREATE INDEX IF NOT EXISTS idx_users_is_teacher ON users(is_teacher);
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);

-- =====================================================
-- Create Sessions Table
-- =====================================================
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);

-- =====================================================
-- Update Artifacts Table
-- =====================================================

-- First, update existing artifacts to have valid enum values
UPDATE artifacts SET type = 'OTHER' WHERE type NOT IN (
  'INTENTS', 'PLAN', 'EVAL', 'EVAL_ITEM', 'PROMPT', 'PROMPT_POLICY', 
  'QUALITY_METRICS', 'HINT', 'SUMMARY', 'OUTLINE', 'OTHER'
);

-- Change type column to use enum
ALTER TABLE artifacts ALTER COLUMN type TYPE artifact_type USING type::artifact_type;

-- =====================================================
-- Update EvalItemReviews Table  
-- =====================================================

-- First, update existing reviews to have valid enum values
UPDATE eval_item_reviews SET verdict = 'APPROVED' WHERE verdict NOT IN (
  'APPROVED', 'REJECTED', 'NEEDS_REVISION'
);

-- Change verdict column to use enum
ALTER TABLE eval_item_reviews ALTER COLUMN verdict TYPE review_verdict USING verdict::review_verdict;

-- =====================================================
-- Add Comments for Documentation
-- =====================================================

COMMENT ON TYPE artifact_type IS 'Types of artifacts that can be generated in the system';
COMMENT ON TYPE review_verdict IS 'Possible verdicts for eval item reviews';

COMMENT ON TABLE sessions IS 'User authentication sessions';
COMMENT ON COLUMN sessions.token IS 'Unique session token for authentication';
COMMENT ON COLUMN sessions.expires_at IS 'When the session expires';

COMMENT ON COLUMN users.is_admin IS 'True if the user has administrative privileges';
COMMENT ON COLUMN users.is_teacher IS 'True if the user is a teacher/instructor';
COMMENT ON COLUMN users.is_learner IS 'True if the user is a learner/student';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove added columns and tables
ALTER TABLE users DROP COLUMN IF EXISTS is_teacher;
ALTER TABLE users DROP COLUMN IF EXISTS is_learner;

DROP TABLE IF EXISTS sessions;

-- Revert artifacts type to text
ALTER TABLE artifacts ALTER COLUMN type TYPE TEXT;

-- Revert eval_item_reviews verdict to text  
ALTER TABLE eval_item_reviews ALTER COLUMN verdict TYPE TEXT;

-- Drop enums
DROP TYPE IF EXISTS review_verdict;
DROP TYPE IF EXISTS artifact_type;

-- +goose StatementEnd

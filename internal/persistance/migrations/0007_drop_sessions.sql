-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Drop Sessions Table and Related Objects
-- =====================================================

-- Drop triggers first
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_user;
DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_sessions_expires;

-- Drop the sessions table
DROP TABLE IF EXISTS sessions;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- =====================================================
-- Recreate Sessions Table (for rollback)
-- =====================================================

CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

COMMENT ON TABLE sessions IS 'User authentication sessions';
COMMENT ON COLUMN sessions.token IS 'Unique session token for authentication';
COMMENT ON COLUMN sessions.expires_at IS 'When the session expires';

-- +goose StatementEnd

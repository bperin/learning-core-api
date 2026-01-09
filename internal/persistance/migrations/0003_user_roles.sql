-- +goose Up
-- +goose StatementBegin

-- Add role column to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'learner';

-- Add check constraint for allowed roles
-- SQLite doesn't support ALTER TABLE ADD CONSTRAINT easily, but for Postgres:
-- DO $$ BEGIN
--   ALTER TABLE users ADD CONSTRAINT check_user_role CHECK (role IN ('admin', 'teacher', 'learner'));
-- EXCEPTION
--   WHEN others THEN NULL;
-- END $$;

-- Create an index on roles
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

COMMENT ON COLUMN users.role IS 'User role: admin, teacher, or learner';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS role;
-- +goose StatementEnd

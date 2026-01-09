-- +goose Up
-- +goose StatementBegin

-- Remove previous role column if it exists
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Add boolean role flags
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_learner BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_teacher BOOLEAN NOT NULL DEFAULT false;

-- is_admin already exists in 0001_ts_schema.sql, but ensuring consistency
-- ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT false;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_is_learner ON users(is_learner);
CREATE INDEX IF NOT EXISTS idx_users_is_teacher ON users(is_teacher);
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);

COMMENT ON COLUMN users.is_admin IS 'True if the user has administrative privileges';
COMMENT ON COLUMN users.is_teacher IS 'True if the user is a teacher/instructor';
COMMENT ON COLUMN users.is_learner IS 'True if the user is a learner/student';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS is_teacher;
ALTER TABLE users DROP COLUMN IF EXISTS is_learner;
-- +goose StatementEnd

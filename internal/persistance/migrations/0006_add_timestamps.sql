-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Add timestamps to eval_items table
-- =====================================================
ALTER TABLE eval_items ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE eval_items ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_eval_items_created_at ON eval_items(created_at);
CREATE INDEX IF NOT EXISTS idx_eval_items_updated_at ON eval_items(updated_at);

-- =====================================================
-- Add timestamps to user_answers table
-- =====================================================
-- user_answers already has created_at, just add updated_at
ALTER TABLE user_answers ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_user_answers_updated_at ON user_answers(updated_at);

-- =====================================================
-- Add timestamps to eval_item_reviews table
-- =====================================================
-- eval_item_reviews already has created_at, just add updated_at
ALTER TABLE eval_item_reviews ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_eval_item_reviews_updated_at ON eval_item_reviews(updated_at);

-- =====================================================
-- Add timestamps to artifacts table
-- =====================================================
-- artifacts already has created_at, just add updated_at
ALTER TABLE artifacts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_artifacts_updated_at ON artifacts(updated_at);

-- =====================================================
-- Create trigger function for automatic updated_at
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- =====================================================
-- Create triggers for automatic updated_at updates
-- =====================================================

-- Trigger for users table (already exists, but ensure it's there)
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for subjects table (already exists, but ensure it's there)
DROP TRIGGER IF EXISTS update_subjects_updated_at ON subjects;
CREATE TRIGGER update_subjects_updated_at
    BEFORE UPDATE ON subjects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for documents table (already exists, but ensure it's there)
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
CREATE TRIGGER update_documents_updated_at
    BEFORE UPDATE ON documents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for evals table (already exists, but ensure it's there)
DROP TRIGGER IF EXISTS update_evals_updated_at ON evals;
CREATE TRIGGER update_evals_updated_at
    BEFORE UPDATE ON evals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for eval_items table (new)
DROP TRIGGER IF EXISTS update_eval_items_updated_at ON eval_items;
CREATE TRIGGER update_eval_items_updated_at
    BEFORE UPDATE ON eval_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for prompt_templates table (already exists, but ensure it's there)
DROP TRIGGER IF EXISTS update_prompt_templates_updated_at ON prompt_templates;
CREATE TRIGGER update_prompt_templates_updated_at
    BEFORE UPDATE ON prompt_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for test_attempts table (new)
DROP TRIGGER IF EXISTS update_test_attempts_updated_at ON test_attempts;
CREATE TRIGGER update_test_attempts_updated_at
    BEFORE UPDATE ON test_attempts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for user_answers table (new)
DROP TRIGGER IF EXISTS update_user_answers_updated_at ON user_answers;
CREATE TRIGGER update_user_answers_updated_at
    BEFORE UPDATE ON user_answers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for artifacts table (new)
DROP TRIGGER IF EXISTS update_artifacts_updated_at ON artifacts;
CREATE TRIGGER update_artifacts_updated_at
    BEFORE UPDATE ON artifacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for eval_item_reviews table (new)
DROP TRIGGER IF EXISTS update_eval_item_reviews_updated_at ON eval_item_reviews;
CREATE TRIGGER update_eval_item_reviews_updated_at
    BEFORE UPDATE ON eval_item_reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for sessions table (new)
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- Add comments for documentation
-- =====================================================
COMMENT ON COLUMN eval_items.created_at IS 'When the eval item was created';
COMMENT ON COLUMN eval_items.updated_at IS 'When the eval item was last updated';

COMMENT ON COLUMN user_answers.updated_at IS 'When the user answer was last updated';

COMMENT ON COLUMN eval_item_reviews.updated_at IS 'When the review was last updated';

COMMENT ON COLUMN artifacts.updated_at IS 'When the artifact was last updated';

COMMENT ON FUNCTION update_updated_at_column() IS 'Trigger function to automatically update updated_at timestamp';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP TRIGGER IF EXISTS update_eval_item_reviews_updated_at ON eval_item_reviews;
DROP TRIGGER IF EXISTS update_artifacts_updated_at ON artifacts;
DROP TRIGGER IF EXISTS update_user_answers_updated_at ON user_answers;
DROP TRIGGER IF EXISTS update_test_attempts_updated_at ON test_attempts;
DROP TRIGGER IF EXISTS update_eval_items_updated_at ON eval_items;

-- Remove added columns
ALTER TABLE eval_items DROP COLUMN IF EXISTS updated_at;
ALTER TABLE eval_items DROP COLUMN IF EXISTS created_at;

ALTER TABLE user_answers DROP COLUMN IF EXISTS updated_at;
ALTER TABLE eval_item_reviews DROP COLUMN IF EXISTS updated_at;
ALTER TABLE artifacts DROP COLUMN IF EXISTS updated_at;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- +goose StatementEnd

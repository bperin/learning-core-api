-- +goose Up
-- Create eval_prompts table to store versioned evaluation prompts
CREATE TABLE eval_prompts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    eval_type TEXT NOT NULL, -- 'groundedness', 'answerability', 'hierarchy'
    version INT NOT NULL DEFAULT 1,
    prompt_text TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(eval_type, version)
);

-- Create eval_results table to store evaluation outcomes
CREATE TABLE eval_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    eval_item_id UUID NOT NULL REFERENCES eval_items(id),
    eval_type TEXT NOT NULL, -- 'groundedness', 'answerability', 'hierarchy'
    eval_prompt_id UUID NOT NULL REFERENCES eval_prompts(id),
    score FLOAT,
    is_grounded BOOLEAN,
    verdict TEXT, -- 'PASS', 'FAIL', 'WARN'
    reasoning TEXT,
    unsupported_claims JSONB,
    gcp_eval_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX idx_eval_prompts_eval_type_active ON eval_prompts(eval_type, is_active);
CREATE INDEX idx_eval_results_eval_item_id ON eval_results(eval_item_id);
CREATE INDEX idx_eval_results_eval_type ON eval_results(eval_type);

-- Add comments
COMMENT ON TABLE eval_prompts IS 'Versioned evaluation prompts for different eval types';
COMMENT ON TABLE eval_results IS 'Results from running evaluations on eval items';
COMMENT ON COLUMN eval_prompts.eval_type IS 'Type of evaluation: groundedness, answerability, hierarchy';
COMMENT ON COLUMN eval_prompts.is_active IS 'Whether this version is the active/default version';
COMMENT ON COLUMN eval_results.is_grounded IS 'Whether the response is grounded in the context (for groundedness evals)';
COMMENT ON COLUMN eval_results.unsupported_claims IS 'List of claims not supported by context (for groundedness evals)';

-- +goose Down
DROP TABLE IF EXISTS eval_results;
DROP TABLE IF EXISTS eval_prompts;

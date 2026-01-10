-- +goose Up
-- +goose StatementBegin

-- Add mime_type column
ALTER TABLE model_configs ADD COLUMN mime_type TEXT;

-- Seed initial model config
INSERT INTO model_configs (
    key, 
    version, 
    provider, 
    model_name, 
    temperature, 
    max_tokens, 
    top_p, 
    top_k, 
    mime_type,
    is_active, 
    created_by
) VALUES (
    'default_gemini',
    1,
    'gemini',
    'gemini-3-preview',
    1.0,
    8192,
    0.5,
    20,
    'application/json',
    true,
    (SELECT id FROM users LIMIT 1)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM model_configs WHERE key = 'default_gemini';
ALTER TABLE model_configs DROP COLUMN mime_type;
-- +goose StatementEnd
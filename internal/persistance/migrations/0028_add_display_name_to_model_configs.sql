-- Add display_name to model_configs
ALTER TABLE model_configs ADD COLUMN display_name TEXT;

-- Update existing rows to have a default display_name (using model_name)
UPDATE model_configs SET display_name = model_name WHERE display_name IS NULL;

-- Make it NOT NULL for future entries
ALTER TABLE model_configs ALTER COLUMN display_name SET NOT NULL;

-- Set current active model name to empty string (as per feedback)
UPDATE model_configs SET model_name = '' WHERE is_active = true;

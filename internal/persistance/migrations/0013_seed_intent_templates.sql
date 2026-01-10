-- +goose Up
-- +goose StatementBegin

-- Create a system user if not exists (for seeding purposes)
INSERT INTO users (email, password, is_admin)
VALUES ('system@slap.events', 'system_placeholder_password', true)
ON CONFLICT (email) DO NOTHING;

INSERT INTO prompt_templates (key, version, is_active, title, description, template, created_by)
VALUES (
    'intent_extraction', 
    1, 
    true, 
    'Document Intent Extraction', 
    'Analyzes a document to extract summary, topics, and intent.', 
    'Analyze the following document content and extract its intent, summary, and educational context.

Document Content:
{{.DocumentContent}}', 
    'system'
)
ON CONFLICT (key, version) DO NOTHING;

INSERT INTO schema_templates (schema_type, version, schema_json, is_active, created_by)
VALUES (
    'intent_extraction',
    1,
    '{
  "type": "OBJECT",
  "properties": {
    "title": { "type": "STRING", "description": "A concise title for the document content" },
    "summary": { "type": "STRING", "description": "A 2-3 sentence summary of the main points" },
    "key_topics": { 
      "type": "ARRAY", 
      "items": { "type": "STRING" }, 
      "description": "List of 3-5 key topics covered" 
    },
    "intended_audience": { "type": "STRING", "description": "Who is this document for?" },
    "educational_value": { "type": "STRING", "description": "What can be learned from this?" },
    "difficulty_level": { 
        "type": "STRING", 
        "enum": ["BEGINNER", "INTERMEDIATE", "ADVANCED"],
        "description": "Estimated difficulty level"
    }
  },
  "required": ["title", "summary", "key_topics", "intended_audience"]
}',
    true,
    (SELECT id FROM users WHERE email = 'system@slap.events' LIMIT 1)
)
ON CONFLICT (schema_type, version) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM prompt_templates WHERE key = 'intent_extraction';
DELETE FROM schema_templates WHERE schema_type = 'intent_extraction';
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

-- Create a system user if not exists
INSERT INTO users (email, password, is_admin)
VALUES ('system@slap.events', 'system_placeholder_password', true)
ON CONFLICT (email) DO NOTHING;

-- Deactivate previous version of intent_extraction prompt
UPDATE prompt_templates SET is_active = false WHERE key = 'intent_extraction';

-- Insert new version of intent_extraction prompt (RAG based)
INSERT INTO prompt_templates (key, version, is_active, title, description, template, created_by)
VALUES (
    'intent_extraction', 
    2, 
    true, 
    'Educational Intent Extraction (RAG)', 
    'Extracts learning intents from documents via file search/RAG.', 
    'You are an educational content analyst.

You are given one or more source documents via file search.
Your task is to extract learning intents strictly from the content.

Rules:
- Use ONLY the provided documents.
- Do NOT infer beyond the text.
- If information is missing, mark it explicitly as "unknown".
- Do not explain your reasoning.
- Output valid JSON only.', 
    'system'
)
ON CONFLICT (key, version) DO NOTHING;

-- Deactivate previous version of intent_extraction schema
UPDATE schema_templates SET is_active = false WHERE schema_type = 'intent_extraction';

-- Insert new version of intent_extraction schema
INSERT INTO schema_templates (schema_type, version, schema_json, is_active, created_by)
VALUES (
    'intent_extraction',
    2,
    '{
  "type": "OBJECT",
  "properties": {
    "domain": { "type": "STRING", "description": "The broad domain of the content" },
    "subject": { "type": "STRING", "description": "The specific subject area" },
    "intended_audience": { "type": "STRING", "description": "Who the content is designed for" },
    "assumed_prerequisites": { 
      "type": "ARRAY", 
      "items": { "type": "STRING" },
      "description": "What the learner should know before starting"
    },
    "learning_objectives": { 
      "type": "ARRAY", 
      "items": { "type": "STRING" },
      "description": "What the learner will achieve"
    },
    "key_concepts": { 
      "type": "ARRAY", 
      "items": { "type": "STRING" },
      "description": "Fundamental concepts covered in the text"
    },
    "difficulty_level": { 
      "type": "STRING", 
      "enum": ["introductory", "intermediate", "advanced"],
      "description": "The depth/complexity of the material"
    },
    "recommended_artifacts": {
      "type": "OBJECT",
      "properties": {
        "flashcards": { "type": "INTEGER", "description": "Number of flashcards to generate" },
        "multiple_choice_questions": { "type": "INTEGER", "description": "Number of MCQs to generate" },
        "short_answer_questions": { "type": "INTEGER", "description": "Number of short answer questions to generate" }
      },
      "required": ["flashcards", "multiple_choice_questions", "short_answer_questions"]
    }
  },
  "required": ["domain", "subject", "intended_audience", "learning_objectives", "key_concepts", "difficulty_level", "recommended_artifacts"]
}',
    true,
    (SELECT id FROM users WHERE email = 'system@slap.events' LIMIT 1)
)
ON CONFLICT (schema_type, version) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Re-activate version 1 if needed, or just delete version 2
DELETE FROM prompt_templates WHERE key = 'intent_extraction' AND version = 2;
DELETE FROM schema_templates WHERE schema_type = 'intent_extraction' AND version = 2;
UPDATE prompt_templates SET is_active = true WHERE key = 'intent_extraction' AND version = 1;
UPDATE schema_templates SET is_active = true WHERE schema_type = 'intent_extraction' AND version = 1;
-- +goose StatementEnd

package generation

import (
	"encoding/json"

	"github.com/google/uuid"
)

// GenerateRequest is the unified entry point for all AI generation tasks.
// It can either reference existing templates in the database or provide
// explicit inline configuration for instructions, schemas, and tools.
type GenerateRequest struct {
	UserID uuid.UUID `json:"user_id"`

	// Targeting Metadata - identifies what the generation is for
	Target Target `json:"target"`

	// Instructions - either a reference to a template or explicit text
	Instructions Instructions `json:"instructions"`

	// OutputConfig - defines the expected output format (Text vs JSON Schema)
	Output OutputConfig `json:"output"`

	// Tools - optional tools like RAG (File Search) or Function Calling
	Tools []ToolConfig `json:"tools,omitempty"`

	// Model Configuration
	ModelConfigID uuid.UUID `json:"model_config_id"`
}

type Target struct {
	DocumentID *uuid.UUID `json:"document_id,omitempty"`
	EvalID     *uuid.UUID `json:"eval_id,omitempty"`
	EvalItemID *uuid.UUID `json:"eval_item_id,omitempty"`
	AttemptID  *uuid.UUID `json:"attempt_id,omitempty"`
}

type Instructions struct {
	SystemInstructionID *uuid.UUID             `json:"system_instruction_id,omitempty"`
	GenerationType      string                 `json:"generation_type,omitempty"` // Reference to DB template
	PromptVersion       int32                  `json:"prompt_version,omitempty"`  // 0 for latest
	Variables           map[string]interface{} `json:"variables,omitempty"`       // Variables to inject into template
	Inline              string                 `json:"inline,omitempty"`          // Raw prompt text (if not using generation type)
}

type OutputConfig struct {
	GenerationType string          `json:"generation_type,omitempty"` // Reference to DB schema
	SchemaVersion  int32           `json:"schema_version,omitempty"`  // 0 for latest
	InlineSchema   json.RawMessage `json:"inline_schema,omitempty"`   // Raw JSON Schema
	Format         string          `json:"format"`                    // "text" or "json"
}

type ToolConfig struct {
	Type   string          `json:"type"`             // "file_search", "function_calling"
	Config json.RawMessage `json:"config,omitempty"` // Tool-specific configuration
}

type ModelConfig struct {
	Name        string   `json:"name,omitempty"` // e.g., "gemini-1.5-pro"
	Temperature *float32 `json:"temperature,omitempty"`
	MaxTokens   *int32   `json:"max_tokens,omitempty"`
	TopP        *float32 `json:"top_p,omitempty"`
	TopK        *float32 `json:"top_k,omitempty"`
	MimeType    string   `json:"mime_type,omitempty"` // e.g. "application/json"
}

type GenerateResponse struct {
	ArtifactID   uuid.UUID       `json:"artifact_id"`
	OutputText   string          `json:"output_text"`
	OutputJSON   json.RawMessage `json:"output_json,omitempty"`
	FinishReason string          `json:"finish_reason"`
	ModelUsed    string          `json:"model_used"`
}

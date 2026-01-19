package schema_templates

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SchemaTemplate represents a versioned schema definition for AI outputs.
// @Description Schema template with versioning and activation support
type SchemaTemplate struct {
	ID             uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GenerationType string          `json:"generation_type" example:"CLASSIFICATION"`
	Version        int32           `json:"version" example:"1"`
	SchemaJSON     json.RawMessage `json:"schema_json" swaggertype:"object"`
	IsActive       bool            `json:"is_active" example:"true"`
	CreatedBy      uuid.UUID       `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt      time.Time       `json:"created_at" example:"2026-01-19T03:40:00Z"`
	LockedAt       *time.Time      `json:"locked_at,omitempty" example:"2026-01-19T03:40:00Z"`
}

// CreateSchemaTemplateRequest represents the data needed to create a schema template.
type CreateSchemaTemplateRequest struct {
	GenerationType string          `json:"generation_type"`
	SchemaJSON     json.RawMessage `json:"schema_json" swaggertype:"object"`
	IsActive       *bool           `json:"is_active,omitempty"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	LockedAt       *time.Time      `json:"locked_at,omitempty"`
}

package prompt_templates

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PromptTemplate represents a prompt template stored in the database.
// @Description Prompt template with versioning and activation support
type PromptTemplate struct {
	ID             uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GenerationType string          `json:"generation_type" example:"CLASSIFICATION"`
	Version        int32           `json:"version" example:"1"`
	IsActive       bool            `json:"is_active" example:"true"`
	Title          string          `json:"title" example:"Classification Prompt v1"`
	Description    *string         `json:"description,omitempty" example:"Classifies documents"`
	Template       string          `json:"template" example:"Classify the following text..."`
	Metadata       json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedBy      *string         `json:"created_by,omitempty" example:"admin@example.com"`
	CreatedAt      time.Time       `json:"created_at" example:"2026-01-19T03:40:00Z"`
	UpdatedAt      time.Time       `json:"updated_at" example:"2026-01-19T03:40:00Z"`
}

// CreatePromptTemplateRequest represents data needed to create a prompt template.
type CreatePromptTemplateRequest struct {
	GenerationType string          `json:"generation_type"`
	Version        int32           `json:"version"`
	IsActive       bool            `json:"is_active"`
	Title          string          `json:"title"`
	Description    *string         `json:"description,omitempty"`
	Template       string          `json:"template"`
	Metadata       json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedBy      *string         `json:"created_by,omitempty"`
}

// CreatePromptTemplateVersionRequest represents data needed to create a new version.
type CreatePromptTemplateVersionRequest struct {
	GenerationType string          `json:"generation_type"`
	IsActive       bool            `json:"is_active"`
	Title          string          `json:"title"`
	Description    *string         `json:"description,omitempty"`
	Template       string          `json:"template"`
	Metadata       json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedBy      *string         `json:"created_by,omitempty"`
}

package prompt_templates

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PromptTemplate represents a prompt template stored in the database.
type PromptTemplate struct {
	ID          uuid.UUID       `json:"id"`
	Key         string          `json:"key"`
	Version     int32           `json:"version"`
	IsActive    bool            `json:"is_active"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Template    string          `json:"template"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedBy   *string         `json:"created_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CreatePromptTemplateRequest represents data needed to create a prompt template.
type CreatePromptTemplateRequest struct {
	Key         string          `json:"key"`
	Version     int32           `json:"version"`
	IsActive    bool            `json:"is_active"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Template    string          `json:"template"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedBy   *string         `json:"created_by,omitempty"`
}

// CreatePromptTemplateVersionRequest represents data needed to create a new version.
type CreatePromptTemplateVersionRequest struct {
	Key         string          `json:"key"`
	IsActive    bool            `json:"is_active"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Template    string          `json:"template"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedBy   *string         `json:"created_by,omitempty"`
}

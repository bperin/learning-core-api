package schema_templates

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SchemaTemplate represents a versioned schema definition for AI outputs.
type SchemaTemplate struct {
	ID           uuid.UUID       `json:"id"`
	SchemaType   string          `json:"schema_type"`
	Version      int32           `json:"version"`
	SchemaJSON   json.RawMessage `json:"schema_json"`
	SubjectID    *uuid.UUID      `json:"subject_id,omitempty"`
	CurriculumID *uuid.UUID      `json:"curriculum_id,omitempty"`
	IsActive     bool            `json:"is_active"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	LockedAt     *time.Time      `json:"locked_at,omitempty"`
}

// CreateSchemaTemplateRequest represents the data needed to create a schema template.
type CreateSchemaTemplateRequest struct {
	SchemaType   string          `json:"schema_type"`
	Version      int32           `json:"version"`
	SchemaJSON   json.RawMessage `json:"schema_json"`
	SubjectID    *uuid.UUID      `json:"subject_id,omitempty"`
	CurriculumID *uuid.UUID      `json:"curriculum_id,omitempty"`
	IsActive     bool            `json:"is_active"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	LockedAt     *time.Time      `json:"locked_at,omitempty"`
}

// CreateSchemaTemplateVersionRequest represents the data needed to create a new version.
type CreateSchemaTemplateVersionRequest struct {
	SchemaType   string          `json:"schema_type"`
	SchemaJSON   json.RawMessage `json:"schema_json"`
	SubjectID    *uuid.UUID      `json:"subject_id,omitempty"`
	CurriculumID *uuid.UUID      `json:"curriculum_id,omitempty"`
	IsActive     bool            `json:"is_active"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	LockedAt     *time.Time      `json:"locked_at,omitempty"`
}

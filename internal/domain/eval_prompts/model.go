package eval_prompts

import (
	"time"

	"github.com/google/uuid"
)

// EvalPrompt represents a versioned evaluation prompt template
type EvalPrompt struct {
	ID          uuid.UUID  `json:"id"`
	EvalType    string     `json:"eval_type"`
	Version     int32      `json:"version"`
	PromptText  string     `json:"prompt_text"`
	Description *string    `json:"description,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateEvalPromptRequest represents a request to create a new eval prompt
type CreateEvalPromptRequest struct {
	EvalType    string     `json:"eval_type" validate:"required,min=1"`
	PromptText  string     `json:"prompt_text" validate:"required,min=1"`
	Description *string    `json:"description,omitempty"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
}

// Validate validates the CreateEvalPromptRequest
func (r *CreateEvalPromptRequest) Validate() error {
	if r.EvalType == "" {
		return ErrInvalidEvalType
	}
	if r.PromptText == "" {
		return ErrEmptyPromptText
	}
	return nil
}

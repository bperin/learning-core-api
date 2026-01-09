package eval_items

import (
	"time"

	"github.com/google/uuid"
)

// EvalItem represents a single question or prompt within an evaluation
type EvalItem struct {
	ID          uuid.UUID `json:"id"`
	EvalID      uuid.UUID `json:"eval_id"`
	Prompt      string    `json:"prompt"`
	Options     []string  `json:"options"`
	CorrectIdx  int32     `json:"correct_idx"`
	Hint        *string   `json:"hint,omitempty"`
	Explanation *string   `json:"explanation,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateEvalItemRequest represents the data needed to create a new evaluation item
type CreateEvalItemRequest struct {
	EvalID      uuid.UUID              `json:"eval_id" validate:"required"`
	Prompt      string                 `json:"prompt" validate:"required,min=1,max=2000"`
	Options     []string               `json:"options" validate:"required,min=2,max=10"`
	CorrectIdx  int32                  `json:"correct_idx" validate:"required,min=0"`
	Hint        *string                `json:"hint,omitempty" validate:"omitempty,max=500"`
	Explanation *string                `json:"explanation,omitempty" validate:"omitempty,max=1000"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateEvalItemRequest represents the data that can be updated for an evaluation item
type UpdateEvalItemRequest struct {
	Prompt      *string                `json:"prompt,omitempty" validate:"omitempty,min=1,max=2000"`
	Options     []string               `json:"options,omitempty" validate:"omitempty,min=2,max=10"`
	CorrectIdx  *int32                 `json:"correct_idx,omitempty" validate:"omitempty,min=0"`
	Hint        *string                `json:"hint,omitempty" validate:"omitempty,max=500"`
	Explanation *string                `json:"explanation,omitempty" validate:"omitempty,max=1000"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ListEvalItemsRequest represents parameters for listing evaluation items
type ListEvalItemsRequest struct {
	EvalID *uuid.UUID `json:"eval_id,omitempty"`
	Limit  int32      `json:"limit" validate:"min=1,max=100"`
	Offset int32      `json:"offset" validate:"min=0"`
}

// SearchEvalItemsRequest represents parameters for searching evaluation items
type SearchEvalItemsRequest struct {
	Query  string `json:"query" validate:"required,min=1"`
	Limit  int32  `json:"limit" validate:"min=1,max=100"`
	Offset int32  `json:"offset" validate:"min=0"`
}

// Validate validates the CreateEvalItemRequest
func (r *CreateEvalItemRequest) Validate() error {
	if r.EvalID == uuid.Nil {
		return ErrInvalidEvalID
	}
	
	if len(r.Prompt) == 0 {
		return ErrEmptyPrompt
	}
	
	if len(r.Options) < 2 {
		return ErrInsufficientOptions
	}
	
	if r.CorrectIdx < 0 || int(r.CorrectIdx) >= len(r.Options) {
		return ErrInvalidCorrectIndex
	}
	
	return nil
}

// Validate validates the UpdateEvalItemRequest
func (r *UpdateEvalItemRequest) Validate() error {
	if r.Prompt != nil && len(*r.Prompt) == 0 {
		return ErrEmptyPrompt
	}
	
	if r.Options != nil && len(r.Options) < 2 {
		return ErrInsufficientOptions
	}
	
	if r.CorrectIdx != nil && r.Options != nil {
		if *r.CorrectIdx < 0 || int(*r.CorrectIdx) >= len(r.Options) {
			return ErrInvalidCorrectIndex
		}
	}
	
	return nil
}

// IsMultipleChoice returns true if the eval item is a multiple choice question
func (e *EvalItem) IsMultipleChoice() bool {
	return len(e.Options) > 0
}

// GetCorrectAnswer returns the correct answer text for multiple choice questions
func (e *EvalItem) GetCorrectAnswer() string {
	if !e.IsMultipleChoice() || int(e.CorrectIdx) >= len(e.Options) {
		return ""
	}
	return e.Options[e.CorrectIdx]
}

// HasHint returns true if the eval item has a hint
func (e *EvalItem) HasHint() bool {
	return e.Hint != nil && *e.Hint != ""
}

// HasExplanation returns true if the eval item has an explanation
func (e *EvalItem) HasExplanation() bool {
	return e.Explanation != nil && *e.Explanation != ""
}

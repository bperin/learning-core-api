package evals

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EvalStatus represents the status of an evaluation
type EvalStatus string

const (
	EvalStatusDraft     EvalStatus = "draft"
	EvalStatusPublished EvalStatus = "published"
	EvalStatusArchived  EvalStatus = "archived"
)

// DifficultyLevel represents the difficulty level of an evaluation
type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "easy"
	DifficultyMedium DifficultyLevel = "medium"
	DifficultyHard   DifficultyLevel = "hard"
)

// Eval represents an evaluation in the domain
type Eval struct {
	ID           uuid.UUID        `json:"id"`
	Title        string           `json:"title"`
	Description  *string          `json:"description,omitempty"`
	Status       EvalStatus       `json:"status"`
	Difficulty   *DifficultyLevel `json:"difficulty,omitempty"`
	Instructions *string          `json:"instructions,omitempty"`
	Rubric       json.RawMessage  `json:"rubric,omitempty"`
	SubjectID    *uuid.UUID       `json:"subject_id,omitempty"`
	UserID       uuid.UUID        `json:"user_id"`
	PublishedAt  *time.Time       `json:"published_at,omitempty"`
	ArchivedAt   *time.Time       `json:"archived_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// EvalWithItemCount represents an evaluation with its item count
type EvalWithItemCount struct {
	*Eval
	ItemCount int64 `json:"item_count"`
}

// CreateEvalRequest represents the request to create an evaluation
type CreateEvalRequest struct {
	Title        string           `json:"title" validate:"required,min=1,max=255"`
	Description  *string          `json:"description,omitempty" validate:"omitempty,max=1000"`
	Difficulty   *DifficultyLevel `json:"difficulty,omitempty"`
	Instructions *string          `json:"instructions,omitempty" validate:"omitempty,max=5000"`
	Rubric       json.RawMessage  `json:"rubric,omitempty"`
	SubjectID    *uuid.UUID       `json:"subject_id,omitempty"`
	UserID       uuid.UUID        `json:"user_id" validate:"required"`
}

// UpdateEvalRequest represents the request to update an evaluation
type UpdateEvalRequest struct {
	Title        *string          `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description  *string          `json:"description,omitempty" validate:"omitempty,max=1000"`
	Difficulty   *DifficultyLevel `json:"difficulty,omitempty"`
	Instructions *string          `json:"instructions,omitempty" validate:"omitempty,max=5000"`
	Rubric       json.RawMessage  `json:"rubric,omitempty"`
	SubjectID    *uuid.UUID       `json:"subject_id,omitempty"`
}

// EvalFilter represents filters for listing evaluations
type EvalFilter struct {
	UserID    *uuid.UUID   `json:"user_id,omitempty"`
	SubjectID *uuid.UUID   `json:"subject_id,omitempty"`
	Status    *EvalStatus  `json:"status,omitempty"`
	Search    *string      `json:"search,omitempty"`
	Limit     int          `json:"limit"`
	Offset    int          `json:"offset"`
}

// IsValidStatus checks if the status is valid
func (s EvalStatus) IsValid() bool {
	switch s {
	case EvalStatusDraft, EvalStatusPublished, EvalStatusArchived:
		return true
	default:
		return false
	}
}

// IsValidDifficulty checks if the difficulty level is valid
func (d DifficultyLevel) IsValid() bool {
	switch d {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the eval can transition to the target status
func (e *Eval) CanTransitionTo(targetStatus EvalStatus) bool {
	switch e.Status {
	case EvalStatusDraft:
		return targetStatus == EvalStatusPublished || targetStatus == EvalStatusArchived
	case EvalStatusPublished:
		return targetStatus == EvalStatusArchived
	case EvalStatusArchived:
		return false // Archived evals cannot transition to any other status
	default:
		return false
	}
}

// CanBeModified checks if the eval can be modified
func (e *Eval) CanBeModified() bool {
	return e.Status == EvalStatusDraft
}

// CanBeDeleted checks if the eval can be deleted
func (e *Eval) CanBeDeleted() bool {
	return e.Status == EvalStatusDraft
}

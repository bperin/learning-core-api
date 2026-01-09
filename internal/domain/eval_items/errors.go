package eval_items

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Domain-specific errors for eval items
var (
	// Validation errors
	ErrInvalidEvalID        = errors.New("invalid evaluation ID")
	ErrEmptyPrompt          = errors.New("prompt cannot be empty")
	ErrInsufficientOptions  = errors.New("at least 2 options are required for multiple choice questions")
	ErrInvalidCorrectIndex  = errors.New("correct index must be within the range of available options")
	ErrPromptTooLong        = errors.New("prompt exceeds maximum length")
	ErrHintTooLong          = errors.New("hint exceeds maximum length")
	ErrExplanationTooLong   = errors.New("explanation exceeds maximum length")
	ErrTooManyOptions       = errors.New("too many options provided")

	// Business logic errors
	ErrEvalItemNotFound     = errors.New("evaluation item not found")
	ErrEvalNotFound         = errors.New("evaluation not found")
	ErrEvalNotDraft         = errors.New("evaluation items can only be modified when evaluation is in draft status")
	ErrCannotDeletePublished = errors.New("cannot delete evaluation items from published evaluations")
	ErrDuplicatePrompt      = errors.New("duplicate prompt within the same evaluation")

	// Permission errors
	ErrUnauthorized         = errors.New("unauthorized to access this evaluation item")
	ErrInsufficientPermissions = errors.New("insufficient permissions to perform this action")
)

// EvalItemNotFoundError represents a specific eval item not found error
type EvalItemNotFoundError struct {
	ID uuid.UUID
}

func (e EvalItemNotFoundError) Error() string {
	return fmt.Sprintf("evaluation item with ID %s not found", e.ID)
}

// NewEvalItemNotFoundError creates a new EvalItemNotFoundError
func NewEvalItemNotFoundError(id uuid.UUID) error {
	return EvalItemNotFoundError{ID: id}
}

// EvalNotFoundError represents a specific evaluation not found error
type EvalNotFoundError struct {
	ID uuid.UUID
}

func (e EvalNotFoundError) Error() string {
	return fmt.Sprintf("evaluation with ID %s not found", e.ID)
}

// NewEvalNotFoundError creates a new EvalNotFoundError
func NewEvalNotFoundError(id uuid.UUID) error {
	return EvalNotFoundError{ID: id}
}

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) error {
	return ValidationError{Field: field, Message: message}
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var evalItemNotFound EvalItemNotFoundError
	var evalNotFound EvalNotFoundError
	return errors.As(err, &evalItemNotFound) || 
		   errors.As(err, &evalNotFound) || 
		   errors.Is(err, ErrEvalItemNotFound) ||
		   errors.Is(err, ErrEvalNotFound)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr ValidationError
	return errors.As(err, &validationErr) ||
		   errors.Is(err, ErrInvalidEvalID) ||
		   errors.Is(err, ErrEmptyPrompt) ||
		   errors.Is(err, ErrInsufficientOptions) ||
		   errors.Is(err, ErrInvalidCorrectIndex) ||
		   errors.Is(err, ErrPromptTooLong) ||
		   errors.Is(err, ErrHintTooLong) ||
		   errors.Is(err, ErrExplanationTooLong) ||
		   errors.Is(err, ErrTooManyOptions)
}

// IsPermissionError checks if an error is a permission error
func IsPermissionError(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		   errors.Is(err, ErrInsufficientPermissions)
}

// IsBusinessLogicError checks if an error is a business logic error
func IsBusinessLogicError(err error) bool {
	return errors.Is(err, ErrEvalNotDraft) ||
		   errors.Is(err, ErrCannotDeletePublished) ||
		   errors.Is(err, ErrDuplicatePrompt)
}

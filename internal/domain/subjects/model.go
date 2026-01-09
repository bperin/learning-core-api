package subjects

import (
	"time"

	"github.com/google/uuid"
)

// Subject represents a subject in the domain
type Subject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSubjectRequest represents the request to create a subject
type CreateSubjectRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
}

// UpdateSubjectRequest represents the request to update a subject
type UpdateSubjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// SubjectFilter represents filters for listing subjects
type SubjectFilter struct {
	UserID uuid.UUID `json:"user_id"`
	Name   *string   `json:"name,omitempty"`
}

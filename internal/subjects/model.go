package subjects

import (
	"time"

	"github.com/google/uuid"
)

// Subject represents a subject scoped to a user.
type Subject struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateSubjectRequest represents the request to create a subject.
type CreateSubjectRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// UpdateSubjectRequest represents the request to update a subject.
type UpdateSubjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

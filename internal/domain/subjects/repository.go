package subjects

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for subject persistence operations
type Repository interface {
	// Create creates a new subject
	Create(ctx context.Context, req CreateSubjectRequest) (*Subject, error)
	
	// GetByID retrieves a subject by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Subject, error)
	
	// Update updates an existing subject
	Update(ctx context.Context, id uuid.UUID, req UpdateSubjectRequest) (*Subject, error)
	
	// Delete deletes a subject by its ID
	Delete(ctx context.Context, id uuid.UUID) error
	
	// ListByUser retrieves all subjects for a specific user
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Subject, error)
}

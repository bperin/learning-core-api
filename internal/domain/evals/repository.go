package evals

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for evaluation persistence operations
type Repository interface {
	// Create creates a new evaluation (always starts as draft)
	Create(ctx context.Context, req CreateEvalRequest) (*Eval, error)

	// GetByID retrieves an evaluation by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Eval, error)

	// GetWithItemCount retrieves an evaluation with its item count
	GetWithItemCount(ctx context.Context, id uuid.UUID) (*EvalWithItemCount, error)

	// Update updates an existing evaluation (only if status is draft)
	Update(ctx context.Context, id uuid.UUID, req UpdateEvalRequest) (*Eval, error)

	// Publish transitions an evaluation from draft to published
	Publish(ctx context.Context, id uuid.UUID) (*Eval, error)

	// Archive transitions an evaluation to archived status
	Archive(ctx context.Context, id uuid.UUID) (*Eval, error)

	// Delete deletes an evaluation (only if status is draft)
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByUser retrieves all evaluations for a specific user
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Eval, error)

	// ListByUserWithItemCounts retrieves evaluations for a user with item counts
	ListByUserWithItemCounts(ctx context.Context, userID uuid.UUID) ([]*EvalWithItemCount, error)

	// ListByStatus retrieves all evaluations with a specific status
	ListByStatus(ctx context.Context, status EvalStatus) ([]*Eval, error)

	// ListPublished retrieves all published evaluations
	ListPublished(ctx context.Context) ([]*Eval, error)

	// ListDrafts retrieves all draft evaluations
	ListDrafts(ctx context.Context) ([]*Eval, error)

	// Search searches evaluations by title
	Search(ctx context.Context, query string, limit, offset int) ([]*Eval, error)

	// List retrieves evaluations with pagination
	List(ctx context.Context, limit, offset int) ([]*Eval, error)
}

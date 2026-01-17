package eval_results

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for eval result data operations
type Repository interface {
	// GetByID retrieves an eval result by ID
	GetByID(ctx context.Context, id uuid.UUID) (*EvalResult, error)

	// GetByEvalItem retrieves all eval results for a specific eval item
	GetByEvalItem(ctx context.Context, evalItemID uuid.UUID) ([]*EvalResult, error)

	// GetLatestByEvalItem retrieves the latest eval result for a specific eval item and type
	GetLatestByEvalItem(ctx context.Context, evalItemID uuid.UUID, evalType string) (*EvalResult, error)

	// ListByType retrieves eval results of a specific type
	ListByType(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalResult, error)

	// List retrieves all eval results
	List(ctx context.Context, limit int32, offset int32) ([]*EvalResult, error)

	// Create creates a new eval result
	Create(ctx context.Context, req *CreateEvalResultRequest) (*EvalResult, error)

	// GetStats retrieves aggregate statistics for eval results of a type
	GetStats(ctx context.Context, evalType string) (*EvalResultStats, error)

	// Count returns the total number of eval results
	Count(ctx context.Context) (int64, error)

	// CountByType returns the count of eval results for a specific type
	CountByType(ctx context.Context, evalType string) (int64, error)
}

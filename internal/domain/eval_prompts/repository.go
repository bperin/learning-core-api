package eval_prompts

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for eval prompt data operations
type Repository interface {
	// GetByID retrieves an eval prompt by ID
	GetByID(ctx context.Context, id uuid.UUID) (*EvalPrompt, error)

	// GetActiveByType retrieves the active eval prompt for a given eval type
	GetActiveByType(ctx context.Context, evalType string) (*EvalPrompt, error)

	// GetByTypeAndVersion retrieves a specific version of an eval prompt
	GetByTypeAndVersion(ctx context.Context, evalType string, version int32) (*EvalPrompt, error)

	// ListByType retrieves all versions of an eval prompt type
	ListByType(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalPrompt, error)

	// Create creates a new eval prompt
	Create(ctx context.Context, req *CreateEvalPromptRequest) (*EvalPrompt, error)

	// Activate activates a specific eval prompt version
	Activate(ctx context.Context, id uuid.UUID) error

	// Deactivate deactivates a specific eval prompt version
	Deactivate(ctx context.Context, id uuid.UUID) error

	// GetLatestVersion gets the latest version number for an eval type
	GetLatestVersion(ctx context.Context, evalType string) (int32, error)
}

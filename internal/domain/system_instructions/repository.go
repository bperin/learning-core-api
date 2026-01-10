package system_instructions

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines system instruction persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreateSystemInstructionRequest) (*SystemInstruction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SystemInstruction, error)
	GetActive(ctx context.Context) (*SystemInstruction, error)
	ListAll(ctx context.Context) ([]*SystemInstruction, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

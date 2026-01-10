package chunking_configs

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines chunking config persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreateChunkingConfigRequest) (*ChunkingConfig, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ChunkingConfig, error)
	GetActive(ctx context.Context) (*ChunkingConfig, error)
	ListAll(ctx context.Context) ([]*ChunkingConfig, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

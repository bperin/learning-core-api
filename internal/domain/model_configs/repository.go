package model_configs

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines model config persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreateModelConfigRequest) (*ModelConfig, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ModelConfig, error)
	GetActive(ctx context.Context) (*ModelConfig, error)
	ListAll(ctx context.Context) ([]*ModelConfig, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

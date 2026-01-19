package chunking_configs

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for chunking configs.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*ChunkingConfig, error)
	GetActive(ctx context.Context) (*ChunkingConfig, error)
	ListAll(ctx context.Context) ([]*ChunkingConfig, error)
	Create(ctx context.Context, req CreateChunkingConfigRequest) (*ChunkingConfig, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

// ServiceImpl implements Service.
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new chunking configs service.
func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// GetByID retrieves a chunking config by ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*ChunkingConfig, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActive retrieves the active chunking config.
func (s *ServiceImpl) GetActive(ctx context.Context) (*ChunkingConfig, error) {
	return s.repo.GetActive(ctx)
}

// ListAll lists all chunking configs.
func (s *ServiceImpl) ListAll(ctx context.Context) ([]*ChunkingConfig, error) {
	return s.repo.ListAll(ctx)
}

// Create creates a new chunking config.
func (s *ServiceImpl) Create(ctx context.Context, req CreateChunkingConfigRequest) (*ChunkingConfig, error) {
	return s.repo.Create(ctx, req)
}

// Activate marks a chunking config as active.
func (s *ServiceImpl) Activate(ctx context.Context, id uuid.UUID) error {
	return s.repo.Activate(ctx, id)
}

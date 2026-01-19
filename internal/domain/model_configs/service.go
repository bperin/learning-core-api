package model_configs

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for model configs.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*ModelConfig, error)
	GetActive(ctx context.Context) (*ModelConfig, error)
	ListAll(ctx context.Context) ([]*ModelConfig, error)
	Create(ctx context.Context, req CreateModelConfigRequest) (*ModelConfig, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

// ServiceImpl implements Service.
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new model configs service.
func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// GetByID retrieves a model config by ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*ModelConfig, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActive retrieves the active model config.
func (s *ServiceImpl) GetActive(ctx context.Context) (*ModelConfig, error) {
	return s.repo.GetActive(ctx)
}

// ListAll lists all model configs.
func (s *ServiceImpl) ListAll(ctx context.Context) ([]*ModelConfig, error) {
	return s.repo.ListAll(ctx)
}

// Create creates a new model config.
func (s *ServiceImpl) Create(ctx context.Context, req CreateModelConfigRequest) (*ModelConfig, error) {
	return s.repo.Create(ctx, req)
}

// Activate marks a model config as active.
func (s *ServiceImpl) Activate(ctx context.Context, id uuid.UUID) error {
	return s.repo.Activate(ctx, id)
}

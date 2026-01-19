package system_instructions

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for system instructions.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*SystemInstruction, error)
	GetActive(ctx context.Context) (*SystemInstruction, error)
	ListAll(ctx context.Context) ([]*SystemInstruction, error)
	Create(ctx context.Context, req CreateSystemInstructionRequest) (*SystemInstruction, error)
	Activate(ctx context.Context, id uuid.UUID) error
}

// ServiceImpl implements Service.
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new system instructions service.
func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// GetByID retrieves a system instruction by ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*SystemInstruction, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActive retrieves the active system instruction.
func (s *ServiceImpl) GetActive(ctx context.Context) (*SystemInstruction, error) {
	return s.repo.GetActive(ctx)
}

// ListAll lists all system instructions.
func (s *ServiceImpl) ListAll(ctx context.Context) ([]*SystemInstruction, error) {
	return s.repo.ListAll(ctx)
}

// Create creates a new system instruction.
func (s *ServiceImpl) Create(ctx context.Context, req CreateSystemInstructionRequest) (*SystemInstruction, error) {
	return s.repo.Create(ctx, req)
}

// Activate marks a system instruction as active.
func (s *ServiceImpl) Activate(ctx context.Context, id uuid.UUID) error {
	return s.repo.Activate(ctx, id)
}

package tenants

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for tenant business logic
type Service interface {
	Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new tenant service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new tenant with validation
func (s *service) Create(ctx context.Context, req CreateTenantRequest) (*Tenant, error) {
	if req.Name == "" {
		return nil, errors.New("tenant name is required")
	}

	tenant := Tenant{
		Name: req.Name,
	}

	return s.repo.Create(ctx, tenant)
}

// GetByID retrieves a tenant by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	if id == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

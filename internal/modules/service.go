package modules

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for module business logic
type Service interface {
	Create(ctx context.Context, req CreateModuleRequest) (*Module, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Module, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*Module, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Module, error)
	Update(ctx context.Context, id uuid.UUID, name, description *string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new module service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new module with business logic validation
func (s *service) Create(ctx context.Context, req CreateModuleRequest) (*Module, error) {
	// Business logic validation
	if req.Name == "" {
		return nil, errors.New("module name is required")
	}

	if req.TenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	// Check if module with same name already exists for this tenant
	existing, err := s.repo.GetByName(ctx, req.TenantID, req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("module with this name already exists for the tenant")
	}

	// Create the module
	module := Module{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Create(ctx, module)
}

// GetByID retrieves a module by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Module, error) {
	if id == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetByName retrieves a module by name and tenant ID
func (s *service) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*Module, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	if name == "" {
		return nil, errors.New("module name is required")
	}

	return s.repo.GetByName(ctx, tenantID, name)
}

// ListByTenant retrieves all modules for a tenant
func (s *service) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Module, error) {
	if tenantID == uuid.Nil {
		return nil, errors.New("tenant ID is required")
	}

	return s.repo.ListByTenant(ctx, tenantID)
}

// Update updates a module with business logic validation
func (s *service) Update(ctx context.Context, id uuid.UUID, name, description *string) error {
	if id == uuid.Nil {
		return errors.New("module ID is required")
	}

	// Check if module exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// If name is being updated, check for conflicts
	if name != nil && *name != existing.Name {
		conflict, err := s.repo.GetByName(ctx, existing.TenantID, *name)
		if err == nil && conflict != nil {
			return errors.New("module with this name already exists for the tenant")
		}
	}

	return s.repo.Update(ctx, id, name, description)
}

// Delete deletes a module
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("module ID is required")
	}

	// Check if module exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

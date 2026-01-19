package schema_templates

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for schema templates.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
	GetActiveByGenerationType(ctx context.Context, generationType string) (*SchemaTemplate, error)
	ListByGenerationType(ctx context.Context, generationType string) ([]*SchemaTemplate, error)
	Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error)
	Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
}

// ServiceImpl implements Service.
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new schema templates service.
func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// GetByID retrieves a schema template by ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActiveByGenerationType retrieves the active schema template for a generation type.
func (s *ServiceImpl) GetActiveByGenerationType(ctx context.Context, generationType string) (*SchemaTemplate, error) {
	return s.repo.GetActiveByGenerationType(ctx, generationType)
}

// ListByGenerationType lists schema templates by generation type.
func (s *ServiceImpl) ListByGenerationType(ctx context.Context, generationType string) ([]*SchemaTemplate, error) {
	return s.repo.ListByGenerationType(ctx, generationType)
}

// Create creates a new schema template.
func (s *ServiceImpl) Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error) {
	return s.repo.Create(ctx, req)
}

// Activate marks a schema template as active.
func (s *ServiceImpl) Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	return s.repo.Activate(ctx, id)
}

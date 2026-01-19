package prompt_templates

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for prompt templates.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
	GetActiveByGenerationType(ctx context.Context, generationType string) (*PromptTemplate, error)
	ListByGenerationType(ctx context.Context, generationType string) ([]*PromptTemplate, error)
	CreateVersion(ctx context.Context, req CreatePromptTemplateVersionRequest) (*PromptTemplate, error)
	Activate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
}

// ServiceImpl implements Service.
type ServiceImpl struct {
	repo Repository
}

// NewService creates a new prompt templates service.
func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// GetByID retrieves a prompt template by ID.
func (s *ServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActiveByGenerationType retrieves the active prompt template for a generation type.
func (s *ServiceImpl) GetActiveByGenerationType(ctx context.Context, generationType string) (*PromptTemplate, error) {
	return s.repo.GetActiveByGenerationType(ctx, generationType)
}

// ListByGenerationType lists all prompt templates for a generation type.
func (s *ServiceImpl) ListByGenerationType(ctx context.Context, generationType string) ([]*PromptTemplate, error) {
	// For now, return the active template for the generation type
	// TODO: Add a ListByGenerationType method to the repository if needed
	template, err := s.repo.GetActiveByGenerationType(ctx, generationType)
	if err != nil {
		return nil, err
	}
	return []*PromptTemplate{template}, nil
}

// CreateVersion creates a new version of a prompt template.
func (s *ServiceImpl) CreateVersion(ctx context.Context, req CreatePromptTemplateVersionRequest) (*PromptTemplate, error) {
	return s.repo.CreateVersion(ctx, req)
}

// Activate marks a prompt template as active.
func (s *ServiceImpl) Activate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	return s.repo.Activate(ctx, id)
}

// Deactivate marks a prompt template as inactive.
func (s *ServiceImpl) Deactivate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	return s.repo.Deactivate(ctx, id)
}

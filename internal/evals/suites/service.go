package suites

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for eval suite business logic
type Service interface {
	Create(ctx context.Context, req CreateSuiteRequest) (*Suite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Suite, error)
	GetByName(ctx context.Context, name string) (*Suite, error)
	List(ctx context.Context) ([]Suite, error)
	Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new eval suite service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new eval suite with business logic validation
func (s *service) Create(ctx context.Context, req CreateSuiteRequest) (*Suite, error) {
	// Business logic validation
	if req.Name == "" {
		return nil, errors.New("suite name is required")
	}

	// Check if eval suite with same name already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("eval suite with this name already exists")
	}

	// Create the suite
	suite := Suite{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Create(ctx, suite)
}

// GetByID retrieves an eval suite by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Suite, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval suite ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetByName retrieves an eval suite by name
func (s *service) GetByName(ctx context.Context, name string) (*Suite, error) {
	if name == "" {
		return nil, errors.New("eval suite name is required")
	}

	return s.repo.GetByName(ctx, name)
}

// List retrieves all eval suites
func (s *service) List(ctx context.Context) ([]Suite, error) {
	return s.repo.List(ctx)
}

// Update updates an eval suite with business logic validation
func (s *service) Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval suite ID is required")
	}

	if name == "" {
		return nil, errors.New("suite name is required")
	}

	// Check if suite exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, name, description)
}

// Delete deletes an eval suite
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("eval suite ID is required")
	}

	// Check if suite exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

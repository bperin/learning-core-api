package evals

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// SuiteService defines the interface for eval suite business logic.
type SuiteService interface {
	Create(ctx context.Context, req CreateSuiteRequest) (*Suite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Suite, error)
	GetByName(ctx context.Context, name string) (*Suite, error)
	List(ctx context.Context) ([]Suite, error)
	Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// suiteService implements the SuiteService interface.
type suiteService struct {
	repo SuiteRepository
}

// NewSuiteService creates a new eval suite service.
func NewSuiteService(repo SuiteRepository) SuiteService {
	return &suiteService{
		repo: repo,
	}
}

// Create creates a new eval suite with business logic validation.
func (s *suiteService) Create(ctx context.Context, req CreateSuiteRequest) (*Suite, error) {
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

// GetByID retrieves an eval suite by ID.
func (s *suiteService) GetByID(ctx context.Context, id uuid.UUID) (*Suite, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval suite ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetByName retrieves an eval suite by name.
func (s *suiteService) GetByName(ctx context.Context, name string) (*Suite, error) {
	if name == "" {
		return nil, errors.New("eval suite name is required")
	}

	return s.repo.GetByName(ctx, name)
}

// List retrieves all eval suites.
func (s *suiteService) List(ctx context.Context) ([]Suite, error) {
	return s.repo.List(ctx)
}

// Update updates an eval suite with business logic validation.
func (s *suiteService) Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error) {
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

// Delete deletes an eval suite.
func (s *suiteService) Delete(ctx context.Context, id uuid.UUID) error {
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

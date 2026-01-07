package evals

import (
	"context"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// SuiteRepository defines the interface for eval suite operations.
type SuiteRepository interface {
	Create(ctx context.Context, suite Suite) (*Suite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Suite, error)
	GetByName(ctx context.Context, name string) (*Suite, error)
	List(ctx context.Context) ([]Suite, error)
	Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// suiteRepository implements the SuiteRepository interface.
type suiteRepository struct {
	queries *store.Queries
}

// NewSuiteRepository creates a new eval suite repository.
func NewSuiteRepository(queries *store.Queries) SuiteRepository {
	return &suiteRepository{
		queries: queries,
	}
}

// Create creates a new eval suite.
func (r *suiteRepository) Create(ctx context.Context, suite Suite) (*Suite, error) {
	params := store.CreateEvalSuiteParams{
		Name:        suite.Name,
		Description: suite.Description,
	}

	dbSuite, err := r.queries.CreateEvalSuite(ctx, params)
	if err != nil {
		return nil, err
	}

	suite := toDomainSuite(dbSuite)
	return &suite, nil
}

// GetByID retrieves an eval suite by ID.
func (r *suiteRepository) GetByID(ctx context.Context, id uuid.UUID) (*Suite, error) {
	dbSuite, err := r.queries.GetEvalSuiteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	suite := toDomainSuite(dbSuite)
	return &suite, nil
}

// GetByName retrieves an eval suite by name.
func (r *suiteRepository) GetByName(ctx context.Context, name string) (*Suite, error) {
	dbSuite, err := r.queries.GetEvalSuiteByName(ctx, name)
	if err != nil {
		return nil, err
	}

	suite := toDomainSuite(dbSuite)
	return &suite, nil
}

// List retrieves all eval suites.
func (r *suiteRepository) List(ctx context.Context) ([]Suite, error) {
	dbSuites, err := r.queries.ListEvalSuites(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainSuites(dbSuites), nil
}

// Update updates an eval suite.
func (r *suiteRepository) Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error) {
	params := store.UpdateEvalSuiteParams{
		ID:          id,
		Name:        name,
		Description: description,
	}

	dbSuite, err := r.queries.UpdateEvalSuite(ctx, params)
	if err != nil {
		return nil, err
	}

	suite := toDomainSuite(dbSuite)
	return &suite, nil
}

// Delete deletes an eval suite.
func (r *suiteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalSuite(ctx, id)
}

func toDomainSuite(dbSuite store.EvalSuite) Suite {
	return Suite{
		ID:          dbSuite.ID,
		Name:        dbSuite.Name,
		Description: dbSuite.Description,
		CreatedAt:   dbSuite.CreatedAt,
	}
}

func toDomainSuites(dbSuites []store.EvalSuite) []Suite {
	suites := make([]Suite, len(dbSuites))
	for i, dbSuite := range dbSuites {
		suites[i] = toDomainSuite(dbSuite)
	}
	return suites
}

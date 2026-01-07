package suites

import (
	"context"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines the interface for eval suite operations
type Repository interface {
	Create(ctx context.Context, suite Suite) (*Suite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Suite, error)
	GetByName(ctx context.Context, name string) (*Suite, error)
	List(ctx context.Context) ([]Suite, error)
	Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new eval suite repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new eval suite
func (r *repository) Create(ctx context.Context, suite Suite) (*Suite, error) {
	params := store.CreateEvalSuiteParams{
		Name:        suite.Name,
		Description: suite.Description,
	}

	dbSuite, err := r.queries.CreateEvalSuite(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Suite{
		ID:          dbSuite.ID,
		Name:        dbSuite.Name,
		Description: dbSuite.Description,
		CreatedAt:   dbSuite.CreatedAt,
	}, nil
}

// GetByID retrieves an eval suite by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Suite, error) {
	dbSuite, err := r.queries.GetEvalSuiteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Suite{
		ID:          dbSuite.ID,
		Name:        dbSuite.Name,
		Description: dbSuite.Description,
		CreatedAt:   dbSuite.CreatedAt,
	}, nil
}

// GetByName retrieves an eval suite by name
func (r *repository) GetByName(ctx context.Context, name string) (*Suite, error) {
	dbSuite, err := r.queries.GetEvalSuiteByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &Suite{
		ID:          dbSuite.ID,
		Name:        dbSuite.Name,
		Description: dbSuite.Description,
		CreatedAt:   dbSuite.CreatedAt,
	}, nil
}

// List retrieves all eval suites
func (r *repository) List(ctx context.Context) ([]Suite, error) {
	dbSuites, err := r.queries.ListEvalSuites(ctx)
	if err != nil {
		return nil, err
	}

	suites := make([]Suite, len(dbSuites))
	for i, dbSuite := range dbSuites {
		suites[i] = Suite{
			ID:          dbSuite.ID,
			Name:        dbSuite.Name,
			Description: dbSuite.Description,
			CreatedAt:   dbSuite.CreatedAt,
		}
	}

	return suites, nil
}

// Update updates an eval suite
func (r *repository) Update(ctx context.Context, id uuid.UUID, name string, description string) (*Suite, error) {
	params := store.UpdateEvalSuiteParams{
		ID:          id,
		Name:        name,
		Description: description,
	}

	dbSuite, err := r.queries.UpdateEvalSuite(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Suite{
		ID:          dbSuite.ID,
		Name:        dbSuite.Name,
		Description: dbSuite.Description,
		CreatedAt:   dbSuite.CreatedAt,
	}, nil
}

// Delete deletes an eval suite
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalSuite(ctx, id)
}

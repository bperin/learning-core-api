package evals

import (
	"context"
	"database/sql"
	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"

	"github.com/google/uuid"
)

// ResultRepository defines the interface for eval result operations.
type ResultRepository interface {
	Create(ctx context.Context, result Result) (*Result, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Result, error)
	ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error)
	ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error)
	DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error
}

// resultRepository implements the ResultRepository interface.
type resultRepository struct {
	queries *store.Queries
}

// NewResultRepository creates a new eval result repository.
func NewResultRepository(queries *store.Queries) ResultRepository {
	return &resultRepository{
		queries: queries,
	}
}

// Create creates a new eval result.
func (r *resultRepository) Create(ctx context.Context, result Result) (*Result, error) {
	params := store.UpsertEvalResultParams{
		EvalRunID: result.EvalRunID,
		RuleID:    result.RuleID,
		Pass:      result.Pass,
		Details:   result.Details,
	}

	if result.Score != nil {
		params.Score = sql.NullFloat64{Float64: float64(*result.Score), Valid: true}
	} else {
		params.Score = sql.NullFloat64{Valid: false}
	}

	dbResult, err := r.queries.UpsertEvalResult(ctx, params)
	if err != nil {
		return nil, err
	}

	result := toDomainResult(dbResult)
	return &result, nil
}

// GetByID retrieves an eval result by ID.
func (r *resultRepository) GetByID(ctx context.Context, id uuid.UUID) (*Result, error) {
	dbResult, err := r.queries.GetEvalResultByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := toDomainResult(dbResult)
	return &result, nil
}

// ListByRun retrieves all eval results for an eval run.
func (r *resultRepository) ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error) {
	dbResults, err := r.queries.ListEvalResultsByRun(ctx, evalRunID)
	if err != nil {
		return nil, err
	}

	return toDomainResults(dbResults), nil
}

// ListByRule retrieves all eval results for a specific rule.
func (r *resultRepository) ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error) {
	dbResults, err := r.queries.ListEvalResultsByRule(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	return toDomainResults(dbResults), nil
}

// DeleteByRun deletes all eval results for an eval run.
func (r *resultRepository) DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error {
	return r.queries.DeleteEvalResultsByRun(ctx, evalRunID)
}

func toDomainResult(dbResult store.EvalResult) Result {
	return Result{
		ID:        dbResult.ID,
		EvalRunID: dbResult.EvalRunID,
		RuleID:    dbResult.RuleID,
		Pass:      dbResult.Pass,
		Score:     utils.NullFloat64ToFloat32Ptr(dbResult.Score),
		Details:   dbResult.Details,
		CreatedAt: dbResult.CreatedAt,
	}
}

func toDomainResults(dbResults []store.EvalResult) []Result {
	results := make([]Result, len(dbResults))
	for i, dbResult := range dbResults {
		results[i] = toDomainResult(dbResult)
	}
	return results
}

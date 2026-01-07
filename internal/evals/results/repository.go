package results

import (
	"context"
	"database/sql"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines the interface for eval result operations
type Repository interface {
	Create(ctx context.Context, result Result) (*Result, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Result, error)
	ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error)
	ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error)
	DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new eval result repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new eval result
func (r *repository) Create(ctx context.Context, result Result) (*Result, error) {
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

	var score *float32
	if dbResult.Score.Valid {
		s := float32(dbResult.Score.Float64)
		score = &s
	}

	return &Result{
		ID:        dbResult.ID,
		EvalRunID: dbResult.EvalRunID,
		RuleID:    dbResult.RuleID,
		Pass:      dbResult.Pass,
		Score:     score,
		Details:   dbResult.Details,
		CreatedAt: dbResult.CreatedAt,
	}, nil
}

// GetByID retrieves an eval result by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Result, error) {
	dbResult, err := r.queries.GetEvalResultByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var score *float32
	if dbResult.Score.Valid {
		s := float32(dbResult.Score.Float64)
		score = &s
	}

	return &Result{
		ID:        dbResult.ID,
		EvalRunID: dbResult.EvalRunID,
		RuleID:    dbResult.RuleID,
		Pass:      dbResult.Pass,
		Score:     score,
		Details:   dbResult.Details,
		CreatedAt: dbResult.CreatedAt,
	}, nil
}

// ListByRun retrieves all eval results for an eval run
func (r *repository) ListByRun(ctx context.Context, evalRunID uuid.UUID) ([]Result, error) {
	dbResults, err := r.queries.ListEvalResultsByRun(ctx, evalRunID)
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(dbResults))
	for i, dbResult := range dbResults {
		var score *float32
		if dbResult.Score.Valid {
			s := float32(dbResult.Score.Float64)
			score = &s
		}

		results[i] = Result{
			ID:        dbResult.ID,
			EvalRunID: dbResult.EvalRunID,
			RuleID:    dbResult.RuleID,
			Pass:      dbResult.Pass,
			Score:     score,
			Details:   dbResult.Details,
			CreatedAt: dbResult.CreatedAt,
		}
	}

	return results, nil
}

// ListByRule retrieves all eval results for a specific rule
func (r *repository) ListByRule(ctx context.Context, ruleID uuid.UUID) ([]Result, error) {
	dbResults, err := r.queries.ListEvalResultsByRule(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(dbResults))
	for i, dbResult := range dbResults {
		var score *float32
		if dbResult.Score.Valid {
			s := float32(dbResult.Score.Float64)
			score = &s
		}

		results[i] = Result{
			ID:        dbResult.ID,
			EvalRunID: dbResult.EvalRunID,
			RuleID:    dbResult.RuleID,
			Pass:      dbResult.Pass,
			Score:     score,
			Details:   dbResult.Details,
			CreatedAt: dbResult.CreatedAt,
		}
	}

	return results, nil
}

// DeleteByRun deletes all eval results for an eval run
func (r *repository) DeleteByRun(ctx context.Context, evalRunID uuid.UUID) error {
	return r.queries.DeleteEvalResultsByRun(ctx, evalRunID)
}

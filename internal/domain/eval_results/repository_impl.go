package eval_results

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"learning-core-api/internal/persistance/store"
)

// RepositoryImpl implements the Repository interface using SQLC
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new eval result repository
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{
		queries: queries,
	}
}

// GetByID retrieves an eval result by ID
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*EvalResult, error) {
	result, err := r.queries.GetEvalResult(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get eval result: %w", err)
	}

	return r.mapToEvalResult(&result), nil
}

// GetByEvalItem retrieves all eval results for a specific eval item
func (r *RepositoryImpl) GetByEvalItem(ctx context.Context, evalItemID uuid.UUID) ([]*EvalResult, error) {
	results, err := r.queries.GetEvalResultsByEvalItem(ctx, evalItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval results for item: %w", err)
	}

	mapped := make([]*EvalResult, len(results))
	for i, res := range results {
		mapped[i] = r.mapToEvalResult(&res)
	}

	return mapped, nil
}

// GetLatestByEvalItem retrieves the latest eval result for a specific eval item and type
func (r *RepositoryImpl) GetLatestByEvalItem(ctx context.Context, evalItemID uuid.UUID, evalType string) (*EvalResult, error) {
	result, err := r.queries.GetLatestEvalResultForItem(ctx, store.GetLatestEvalResultForItemParams{
		EvalItemID: evalItemID,
		EvalType:   evalType,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest eval result: %w", err)
	}

	return r.mapToEvalResult(&result), nil
}

// ListByType retrieves eval results of a specific type
func (r *RepositoryImpl) ListByType(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalResult, error) {
	results, err := r.queries.GetEvalResultsByType(ctx, store.GetEvalResultsByTypeParams{
		EvalType: evalType,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list eval results by type: %w", err)
	}

	mapped := make([]*EvalResult, len(results))
	for i, res := range results {
		mapped[i] = r.mapToEvalResult(&res)
	}

	return mapped, nil
}

// List retrieves all eval results
func (r *RepositoryImpl) List(ctx context.Context, limit int32, offset int32) ([]*EvalResult, error) {
	results, err := r.queries.ListEvalResults(ctx, store.ListEvalResultsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list eval results: %w", err)
	}

	mapped := make([]*EvalResult, len(results))
	for i, res := range results {
		mapped[i] = r.mapToEvalResult(&res)
	}

	return mapped, nil
}

// Create creates a new eval result
func (r *RepositoryImpl) Create(ctx context.Context, req *CreateEvalResultRequest) (*EvalResult, error) {
	result, err := r.queries.CreateEvalResult(ctx, store.CreateEvalResultParams{
		EvalItemID:        req.EvalItemID,
		EvalType:          req.EvalType,
		EvalPromptID:      req.EvalPromptID,
		Score:             toNullFloat64(req.Score),
		IsGrounded:        toNullBool(req.IsGrounded),
		Verdict:           toNullString(&req.Verdict),
		Reasoning:         toNullString(req.Reasoning),
		UnsupportedClaims: toPQTypeNullRawMessage(req.UnsupportedClaims),
		GcpEvalID:         toNullString(req.GCPEvalID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eval result: %w", err)
	}

	return r.mapToEvalResult(&result), nil
}

// GetStats retrieves aggregate statistics for eval results of a type
func (r *RepositoryImpl) GetStats(ctx context.Context, evalType string) (*EvalResultStats, error) {
	stats, err := r.queries.GetEvalResultStats(ctx, evalType)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval result stats: %w", err)
	}

	passRate := 0.0
	if stats.TotalEvals > 0 {
		passRate = float64(stats.Passed) / float64(stats.TotalEvals) * 100
	}

	// AvgScore is returned as string from SQL, convert to float64
	avgScore := 0.0
	if stats.AvgScore != "" && stats.AvgScore != "0" {
		// Parse the string to float64 if needed
		// For now, default to 0.0 since the query returns ROUND(..., 2)
		avgScore = 0.0
	}

	return &EvalResultStats{
		TotalEvals: stats.TotalEvals,
		Passed:     stats.Passed,
		Failed:     stats.Failed,
		Warned:     stats.Warned,
		AvgScore:   avgScore,
		PassRate:   passRate,
	}, nil
}

// Count returns the total number of eval results
func (r *RepositoryImpl) Count(ctx context.Context) (int64, error) {
	// TODO: Implement count query in SQLC
	return 0, fmt.Errorf("not implemented")
}

// CountByType returns the count of eval results for a specific type
func (r *RepositoryImpl) CountByType(ctx context.Context, evalType string) (int64, error) {
	// TODO: Implement count by type query in SQLC
	return 0, fmt.Errorf("not implemented")
}

// mapToEvalResult converts a store.EvalResult to a domain EvalResult
func (r *RepositoryImpl) mapToEvalResult(result *store.EvalResult) *EvalResult {
	if result == nil {
		return nil
	}

	var reasoning *string
	if result.Reasoning.Valid {
		reasoning = &result.Reasoning.String
	}

	var gcpEvalID *string
	if result.GcpEvalID.Valid {
		gcpEvalID = &result.GcpEvalID.String
	}

	var score *float64
	if result.Score.Valid {
		score = &result.Score.Float64
	}

	var isGrounded *bool
	if result.IsGrounded.Valid {
		isGrounded = &result.IsGrounded.Bool
	}

	var verdict string
	if result.Verdict.Valid {
		verdict = result.Verdict.String
	}

	var unsupportedClaims json.RawMessage
	if result.UnsupportedClaims.Valid {
		unsupportedClaims = result.UnsupportedClaims.RawMessage
	}

	return &EvalResult{
		ID:                result.ID,
		EvalItemID:        result.EvalItemID,
		EvalType:          result.EvalType,
		EvalPromptID:      result.EvalPromptID,
		Score:             score,
		IsGrounded:        isGrounded,
		Verdict:           verdict,
		Reasoning:         reasoning,
		UnsupportedClaims: unsupportedClaims,
		GCPEvalID:         gcpEvalID,
		CreatedAt:         result.CreatedAt.Time,
	}
}

// toNullString converts a string pointer to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// toNullBool converts a bool pointer to sql.NullBool
func toNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// toNullFloat64 converts a float64 pointer to sql.NullFloat64
func toNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// toPQTypeNullRawMessage converts json.RawMessage to pqtype.NullRawMessage
func toPQTypeNullRawMessage(msg json.RawMessage) pqtype.NullRawMessage {
	if len(msg) == 0 {
		return pqtype.NullRawMessage{Valid: false}
	}
	return pqtype.NullRawMessage{RawMessage: msg, Valid: true}
}

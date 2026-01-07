package evals

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// RunRepository defines the interface for eval run operations.
type RunRepository interface {
	Create(ctx context.Context, run Run) (*Run, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Run, error)
	GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error)
	ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error)
	UpdateResult(ctx context.Context, id uuid.UUID, status string, overallPass *bool, overallScore *float32, error json.RawMessage) (*Run, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// runRepository implements the RunRepository interface.
type runRepository struct {
	queries *store.Queries
}

// NewRunRepository creates a new eval run repository.
func NewRunRepository(queries *store.Queries) RunRepository {
	return &runRepository{
		queries: queries,
	}
}

// Create creates a new eval run.
func (r *runRepository) Create(ctx context.Context, run Run) (*Run, error) {
	params := store.CreateEvalRunParams{
		ArtifactID:      run.ArtifactID,
		GenerationRunID: utils.PtrToNullUUID(run.GenerationRunID),
		SuiteID:         run.SuiteID,
		JudgeModel:      run.JudgeModel,
		JudgeParams:     run.JudgeParams,
	}

	dbRun, err := r.queries.CreateEvalRun(ctx, params)
	if err != nil {
		return nil, err
	}

	run := toDomainRun(dbRun)
	run.Error = nil // Error is populated when updating results
	return &run, nil
}

// GetByID retrieves an eval run by ID.
func (r *runRepository) GetByID(ctx context.Context, id uuid.UUID) (*Run, error) {
	dbRun, err := r.queries.GetEvalRunByID(ctx, id)
	if err != nil {
		return nil, err
	}

	run := toDomainRun(dbRun)
	return &run, nil
}

// GetLatestForArtifact retrieves the latest eval run for an artifact.
func (r *runRepository) GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error) {
	dbRun, err := r.queries.GetLatestEvalRunForArtifact(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	run := toDomainRun(dbRun)
	return &run, nil
}

// ListByArtifact retrieves all eval runs for an artifact.
func (r *runRepository) ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error) {
	dbRuns, err := r.queries.ListEvalRunsByArtifact(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	return toDomainRuns(dbRuns), nil
}

// UpdateResult updates an eval run's result.
func (r *runRepository) UpdateResult(ctx context.Context, id uuid.UUID, status string, overallPass *bool, overallScore *float32, errorData json.RawMessage) (*Run, error) {
	params := store.UpdateEvalRunResultParams{
		ID:     id,
		Status: status,
	}

	if overallPass != nil {
		params.OverallPass = sql.NullBool{Bool: *overallPass, Valid: true}
	} else {
		params.OverallPass = sql.NullBool{Valid: false}
	}

	if overallScore != nil {
		params.OverallScore = sql.NullFloat64{Float64: float64(*overallScore), Valid: true}
	} else {
		params.OverallScore = sql.NullFloat64{Valid: false}
	}

	params.Error = pqtype.NullRawMessage{RawMessage: errorData, Valid: errorData != nil}

	dbRun, err := r.queries.UpdateEvalRunResult(ctx, params)
	if err != nil {
		return nil, err
	}

	run := toDomainRun(dbRun)
	return &run, nil
}

// Delete deletes an eval run.
func (r *runRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalRun(ctx, id)
}

func toDomainRun(dbRun store.EvalRun) Run {
	return Run{
		ID:              dbRun.ID,
		ArtifactID:      dbRun.ArtifactID,
		GenerationRunID: utils.NullUUIDToPtr(dbRun.GenerationRunID),
		SuiteID:         dbRun.SuiteID,
		JudgeModel:      dbRun.JudgeModel,
		JudgeParams:     dbRun.JudgeParams,
		Status:          dbRun.Status,
		OverallPass:     utils.NullBoolToPtr(dbRun.OverallPass),
		OverallScore:    utils.NullFloat64ToFloat32Ptr(dbRun.OverallScore),
		Error:           dbRun.Error.RawMessage,
		StartedAt:       dbRun.StartedAt.Time,
		FinishedAt:      utils.NullTimeToPtr(dbRun.FinishedAt),
		CreatedAt:       dbRun.CreatedAt,
	}
}

func toDomainRuns(dbRuns []store.EvalRun) []Run {
	runs := make([]Run, len(dbRuns))
	for i, dbRun := range dbRuns {
		runs[i] = toDomainRun(dbRun)
	}
	return runs
}

package runs

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Repository defines the interface for eval run operations
type Repository interface {
	Create(ctx context.Context, run Run) (*Run, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Run, error)
	GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error)
	ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error)
	UpdateResult(ctx context.Context, id uuid.UUID, status string, overallPass *bool, overallScore *float32, error json.RawMessage) (*Run, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new eval run repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new eval run
func (r *repository) Create(ctx context.Context, run Run) (*Run, error) {
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

	return &Run{
		ID:              dbRun.ID,
		ArtifactID:      dbRun.ArtifactID,
		GenerationRunID: utils.NullUUIDToPtr(dbRun.GenerationRunID),
		SuiteID:         dbRun.SuiteID,
		JudgeModel:      dbRun.JudgeModel,
		JudgeParams:     dbRun.JudgeParams,
		Status:          dbRun.Status,
		OverallPass:     utils.NullBoolToPtr(dbRun.OverallPass),
		OverallScore:    utils.NullFloat64ToFloat32Ptr(dbRun.OverallScore),
		Error:           nil, // Error is populated when updating results
		StartedAt:       dbRun.StartedAt.Time,
		FinishedAt:      utils.NullTimeToPtr(dbRun.FinishedAt),
		CreatedAt:       dbRun.CreatedAt,
	}, nil
}

// GetByID retrieves an eval run by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Run, error) {
	dbRun, err := r.queries.GetEvalRunByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Run{
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
	}, nil
}

// GetLatestForArtifact retrieves the latest eval run for an artifact
func (r *repository) GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error) {
	dbRun, err := r.queries.GetLatestEvalRunForArtifact(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	return &Run{
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
	}, nil
}

// ListByArtifact retrieves all eval runs for an artifact
func (r *repository) ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error) {
	dbRuns, err := r.queries.ListEvalRunsByArtifact(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	runs := make([]Run, len(dbRuns))
	for i, dbRun := range dbRuns {
		runs[i] = Run{
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

	return runs, nil
}

// UpdateResult updates an eval run's result
func (r *repository) UpdateResult(ctx context.Context, id uuid.UUID, status string, overallPass *bool, overallScore *float32, errorData json.RawMessage) (*Run, error) {
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

	return &Run{
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
	}, nil
}

// Delete deletes an eval run
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEvalRun(ctx, id)
}

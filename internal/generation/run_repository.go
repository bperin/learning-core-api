package generation

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateGenerationRun creates a new generation run.
func (r *repository) CreateGenerationRun(ctx context.Context, run GenerationRun) (*GenerationRun, error) {
	var promptID uuid.NullUUID
	if run.PromptID != nil {
		promptID = uuid.NullUUID{UUID: *run.PromptID, Valid: true}
	}

	params := store.CreateGenerationRunParams{
		ModuleID:       run.ModuleID,
		AgentName:      run.AgentName,
		AgentVersion:   run.AgentVersion,
		Model:          run.Model,
		ModelParams:    run.ModelParams,
		PromptID:       promptID,
		StoreName:      run.StoreName,
		MetadataFilter: run.MetadataFilter,
		Status:         string(run.Status),
		InputPayload:   run.InputPayload,
	}

	dbRun, err := r.queries.CreateGenerationRun(ctx, params)
	if err != nil {
		return nil, err
	}

	createdRun := convertGenerationRun(dbRun)
	return &createdRun, nil
}

// GetGenerationRunByID retrieves a generation run by ID.
func (r *repository) GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error) {
	dbRun, err := r.queries.GetGenerationRun(ctx, id)
	if err != nil {
		return nil, err
	}

	run := convertGenerationRun(dbRun)
	return &run, nil
}

// ListGenerationRunsByModule retrieves all generation runs for a module.
func (r *repository) ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error) {
	dbRuns, err := r.queries.ListGenerationRunsByModule(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	runs := make([]GenerationRun, len(dbRuns))
	for i, dbRun := range dbRuns {
		runs[i] = convertGenerationRun(dbRun)
	}

	return runs, nil
}

// UpdateGenerationRun updates a generation run.
func (r *repository) UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, errorPayload json.RawMessage, startedAt, finishedAt *time.Time) error {
	// Get current run to preserve values if not updating.
	current, err := r.GetGenerationRunByID(ctx, id)
	if err != nil {
		return err
	}

	params := store.UpdateGenerationRunParams{
		ID:            id,
		Status:        string(current.Status),
		OutputPayload: pqtype.NullRawMessage{RawMessage: current.OutputPayload, Valid: len(current.OutputPayload) > 0},
		Error:         pqtype.NullRawMessage{RawMessage: current.Error, Valid: len(current.Error) > 0},
	}

	if current.StartedAt != nil {
		params.StartedAt = sql.NullTime{Time: *current.StartedAt, Valid: true}
	}
	if current.FinishedAt != nil {
		params.FinishedAt = sql.NullTime{Time: *current.FinishedAt, Valid: true}
	}

	if status != nil {
		params.Status = string(*status)
	}

	if len(outputPayload) > 0 {
		params.OutputPayload = pqtype.NullRawMessage{RawMessage: outputPayload, Valid: true}
	}

	if len(errorPayload) > 0 {
		params.Error = pqtype.NullRawMessage{RawMessage: errorPayload, Valid: true}
	}

	if startedAt != nil {
		params.StartedAt = sql.NullTime{Time: *startedAt, Valid: true}
	}

	if finishedAt != nil {
		params.FinishedAt = sql.NullTime{Time: *finishedAt, Valid: true}
	}

	return r.queries.UpdateGenerationRun(ctx, params)
}

// DeleteGenerationRun deletes a generation run.
func (r *repository) DeleteGenerationRun(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteGenerationRun(ctx, id)
}

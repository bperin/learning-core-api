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

// Repository defines the interface for generation operations
type Repository interface {
	// GenerationRun operations
	CreateGenerationRun(ctx context.Context, run GenerationRun) (*GenerationRun, error)
	GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error)
	ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error)
	UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, errorPayload json.RawMessage, startedAt, finishedAt *time.Time) error
	DeleteGenerationRun(ctx context.Context, id uuid.UUID) error

	// Artifact operations
	CreateArtifact(ctx context.Context, artifact Artifact) (*Artifact, error)
	GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error)
	ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error)
	ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error)
	UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new generation repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Helper functions
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func convertArtifact(dbArtifact store.Artifact) Artifact {
	return Artifact{
		ID:              dbArtifact.ID,
		ModuleID:        dbArtifact.ModuleID,
		GenerationRunID: dbArtifact.GenerationRunID,
		Type:            ArtifactType(dbArtifact.Type),
		Status:          ArtifactStatus(dbArtifact.Status),
		SchemaVersion:   dbArtifact.SchemaVersion,
		Difficulty:      nullStringToPtr(dbArtifact.Difficulty),
		Tags:            dbArtifact.Tags,
		ArtifactPayload: dbArtifact.ArtifactPayload,
		Grounding:       dbArtifact.Grounding,
		EvidenceVersion: nullStringToPtr(dbArtifact.EvidenceVersion),
		ApprovedAt:      &dbArtifact.ApprovedAt.Time,
		RejectedAt:      &dbArtifact.RejectedAt.Time,
		CreatedAt:       dbArtifact.CreatedAt,
	}
}

func convertGenerationRun(dbRun store.GenerationRun) GenerationRun {
	var outputPayload, errorPayload json.RawMessage
	if dbRun.OutputPayload.Valid {
		outputPayload = dbRun.OutputPayload.RawMessage
	}
	if dbRun.Error.Valid {
		errorPayload = dbRun.Error.RawMessage
	}

	var promptID *uuid.UUID
	if dbRun.PromptID.Valid {
		promptID = &dbRun.PromptID.UUID
	}

	var startedAt, finishedAt *time.Time
	if dbRun.StartedAt.Valid {
		startedAt = &dbRun.StartedAt.Time
	}
	if dbRun.FinishedAt.Valid {
		finishedAt = &dbRun.FinishedAt.Time
	}

	return GenerationRun{
		ID:             dbRun.ID,
		ModuleID:       dbRun.ModuleID,
		AgentName:      dbRun.AgentName,
		AgentVersion:   dbRun.AgentVersion,
		Model:          dbRun.Model,
		ModelParams:    dbRun.ModelParams,
		PromptID:       promptID,
		StoreName:      dbRun.StoreName,
		MetadataFilter: dbRun.MetadataFilter,
		Status:         RunStatus(dbRun.Status),
		InputPayload:   dbRun.InputPayload,
		OutputPayload:  outputPayload,
		Error:          errorPayload,
		StartedAt:      startedAt,
		FinishedAt:     finishedAt,
		CreatedAt:      dbRun.CreatedAt,
	}
}

// CreateGenerationRun creates a new generation run
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

// GetGenerationRunByID retrieves a generation run by ID
func (r *repository) GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error) {
	dbRun, err := r.queries.GetGenerationRun(ctx, id)
	if err != nil {
		return nil, err
	}

	run := convertGenerationRun(dbRun)
	return &run, nil
}

// ListGenerationRunsByModule retrieves all generation runs for a module
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

// UpdateGenerationRun updates a generation run
func (r *repository) UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, errorPayload json.RawMessage, startedAt, finishedAt *time.Time) error {
	// Get current run to preserve values if not updating
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

// DeleteGenerationRun deletes a generation run
func (r *repository) DeleteGenerationRun(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteGenerationRun(ctx, id)
}

// CreateArtifact creates a new artifact
func (r *repository) CreateArtifact(ctx context.Context, artifact Artifact) (*Artifact, error) {
	params := store.CreateArtifactParams{
		ModuleID:        artifact.ModuleID,
		GenerationRunID: artifact.GenerationRunID,
		Type:            string(artifact.Type),
		Status:          string(artifact.Status),
		SchemaVersion:   artifact.SchemaVersion,
		Tags:            artifact.Tags,
		ArtifactPayload: artifact.ArtifactPayload,
		Grounding:       artifact.Grounding,
	}

	// Handle nullable fields
	if artifact.Difficulty != nil {
		params.Difficulty = sql.NullString{String: *artifact.Difficulty, Valid: true}
	} else {
		params.Difficulty = sql.NullString{String: "", Valid: false}
	}

	if artifact.EvidenceVersion != nil {
		params.EvidenceVersion = sql.NullString{String: *artifact.EvidenceVersion, Valid: true}
	} else {
		params.EvidenceVersion = sql.NullString{String: "", Valid: false}
	}

	dbArtifact, err := r.queries.CreateArtifact(ctx, params)
	if err != nil {
		return nil, err
	}

	createdArtifact := convertArtifact(dbArtifact)
	return &createdArtifact, nil
}

// GetArtifactByID retrieves an artifact by ID
func (r *repository) GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	dbArtifact, err := r.queries.GetArtifact(ctx, id)
	if err != nil {
		return nil, err
	}

	artifact := convertArtifact(dbArtifact)
	return &artifact, nil
}

// ListArtifactsByModule retrieves all artifacts for a module
func (r *repository) ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error) {
	dbArtifacts, err := r.queries.ListArtifactsByModule(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	artifacts := make([]Artifact, len(dbArtifacts))
	for i, dbArtifact := range dbArtifacts {
		artifacts[i] = convertArtifact(dbArtifact)
	}

	return artifacts, nil
}

// ListArtifactsByModuleAndStatus retrieves all artifacts for a module with a specific status
func (r *repository) ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error) {
	dbArtifacts, err := r.queries.ListArtifactsByModuleAndStatus(ctx, store.ListArtifactsByModuleAndStatusParams{
		ModuleID: moduleID,
		Status:   string(status),
	})
	if err != nil {
		return nil, err
	}

	artifacts := make([]Artifact, len(dbArtifacts))
	for i, dbArtifact := range dbArtifacts {
		artifacts[i] = convertArtifact(dbArtifact)
	}

	return artifacts, nil
}

// UpdateArtifactStatus updates an artifact's status
func (r *repository) UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error {
	params := store.UpdateArtifactStatusParams{
		ID:         id,
		Status:     string(status),
		ApprovedAt: sql.NullTime{Time: time.Time{}, Valid: false}, // Default invalid
		RejectedAt: sql.NullTime{Time: time.Time{}, Valid: false}, // Default invalid
	}

	if approvedAt != nil {
		params.ApprovedAt = sql.NullTime{Time: *approvedAt, Valid: true}
	}
	if rejectedAt != nil {
		params.RejectedAt = sql.NullTime{Time: *rejectedAt, Valid: true}
	}

	return r.queries.UpdateArtifactStatus(ctx, params)
}

// DeleteArtifact deletes an artifact
func (r *repository) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteArtifact(ctx, id)
}

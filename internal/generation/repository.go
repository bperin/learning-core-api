package generation

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"
	"time"

	"github.com/google/uuid"
)

// RunRepository defines storage operations for generation runs.
type RunRepository interface {
	CreateGenerationRun(ctx context.Context, run GenerationRun) (*GenerationRun, error)
	GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error)
	ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error)
	UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, errorPayload json.RawMessage, startedAt, finishedAt *time.Time) error
	DeleteGenerationRun(ctx context.Context, id uuid.UUID) error
}

// ArtifactRepository defines storage operations for artifacts.
type ArtifactRepository interface {
	CreateArtifact(ctx context.Context, artifact Artifact) (*Artifact, error)
	GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error)
	ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error)
	ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error)
	UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
}

// Repository bundles generation run and artifact storage operations.
type Repository interface {
	RunRepository
	ArtifactRepository
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

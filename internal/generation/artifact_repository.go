package generation

import (
	"context"
	"database/sql"
	"learning-core-api/internal/store"
	"time"

	"github.com/google/uuid"
)

// CreateArtifact creates a new artifact.
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

	// Handle nullable fields.
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

// GetArtifactByID retrieves an artifact by ID.
func (r *repository) GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	dbArtifact, err := r.queries.GetArtifact(ctx, id)
	if err != nil {
		return nil, err
	}

	artifact := convertArtifact(dbArtifact)
	return &artifact, nil
}

// ListArtifactsByModule retrieves all artifacts for a module.
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

// ListArtifactsByModuleAndStatus retrieves all artifacts for a module with a specific status.
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

// UpdateArtifactStatus updates an artifact's status.
func (r *repository) UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error {
	params := store.UpdateArtifactStatusParams{
		ID:         id,
		Status:     string(status),
		ApprovedAt: sql.NullTime{Time: time.Time{}, Valid: false},
		RejectedAt: sql.NullTime{Time: time.Time{}, Valid: false},
	}

	if approvedAt != nil {
		params.ApprovedAt = sql.NullTime{Time: *approvedAt, Valid: true}
	}
	if rejectedAt != nil {
		params.RejectedAt = sql.NullTime{Time: *rejectedAt, Valid: true}
	}

	return r.queries.UpdateArtifactStatus(ctx, params)
}

// DeleteArtifact deletes an artifact.
func (r *repository) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteArtifact(ctx, id)
}

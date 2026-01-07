package generation

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CreateArtifact creates a new artifact with business logic validation.
func (s *service) CreateArtifact(ctx context.Context, req CreateArtifactRequest) (*Artifact, error) {
	if req.ModuleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if req.GenerationRunID == uuid.Nil {
		return nil, errors.New("generation run ID is required")
	}

	if req.Type == "" {
		return nil, errors.New("artifact type is required")
	}

	if req.SchemaVersion == "" {
		return nil, errors.New("schema version is required")
	}

	artifact := Artifact{
		ModuleID:        req.ModuleID,
		GenerationRunID: req.GenerationRunID,
		Type:            req.Type,
		Status:          ArtifactStatusPendingEval,
		SchemaVersion:   req.SchemaVersion,
		Difficulty:      req.Difficulty,
		Tags:            req.Tags,
		ArtifactPayload: req.ArtifactPayload,
		Grounding:       req.Grounding,
		EvidenceVersion: req.EvidenceVersion,
		ApprovedAt:      nil,
		RejectedAt:      nil,
	}

	return s.repo.CreateArtifact(ctx, artifact)
}

// GetArtifactByID retrieves an artifact by ID.
func (s *service) GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	if id == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	return s.repo.GetArtifactByID(ctx, id)
}

// ListArtifactsByModule retrieves all artifacts for a module.
func (s *service) ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListArtifactsByModule(ctx, moduleID)
}

// ListArtifactsByModuleAndStatus retrieves all artifacts for a module with a specific status.
func (s *service) ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListArtifactsByModuleAndStatus(ctx, moduleID, status)
}

// UpdateArtifactStatus updates an artifact's status with business logic validation.
func (s *service) UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error {
	if id == uuid.Nil {
		return errors.New("artifact ID is required")
	}

	currentArtifact, err := s.repo.GetArtifactByID(ctx, id)
	if err != nil {
		return err
	}

	if !isArtifactStatusTransitionValid(currentArtifact.Status, status) {
		return errors.New("invalid artifact status transition")
	}

	now := time.Now()
	if status == ArtifactStatusApproved && approvedAt == nil {
		approvedAt = &now
	} else if status == ArtifactStatusRejected && rejectedAt == nil {
		rejectedAt = &now
	}

	return s.repo.UpdateArtifactStatus(ctx, id, status, approvedAt, rejectedAt)
}

// isArtifactStatusTransitionValid checks if an artifact status transition is valid.
func isArtifactStatusTransitionValid(from, to ArtifactStatus) bool {
	validTransitions := map[ArtifactStatus][]ArtifactStatus{
		ArtifactStatusPendingEval: {ArtifactStatusApproved, ArtifactStatusRejected},
		ArtifactStatusApproved:    {},
		ArtifactStatusRejected:    {},
	}

	for _, validTo := range validTransitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// DeleteArtifact deletes an artifact.
func (s *service) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("artifact ID is required")
	}

	_, err := s.repo.GetArtifactByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteArtifact(ctx, id)
}

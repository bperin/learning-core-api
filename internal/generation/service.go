package generation

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for generation business logic
type Service interface {
	// GenerationRun operations
	CreateGenerationRun(ctx context.Context, req CreateGenerationRunRequest) (*GenerationRun, error)
	GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error)
	ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error)
	UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, error json.RawMessage, startedAt, finishedAt *time.Time) error
	DeleteGenerationRun(ctx context.Context, id uuid.UUID) error

	// Artifact operations
	CreateArtifact(ctx context.Context, req CreateArtifactRequest) (*Artifact, error)
	GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error)
	ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error)
	ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error)
	UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new generation service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateGenerationRun creates a new generation run with business logic validation
func (s *service) CreateGenerationRun(ctx context.Context, req CreateGenerationRunRequest) (*GenerationRun, error) {
	// Business logic validation
	if req.ModuleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if req.AgentName == "" {
		return nil, errors.New("agent name is required")
	}

	if req.Model == "" {
		return nil, errors.New("model is required")
	}

	if req.StoreName == "" {
		return nil, errors.New("store name is required")
	}

	// Create the generation run
	run := GenerationRun{
		ModuleID:       req.ModuleID,
		AgentName:      req.AgentName,
		AgentVersion:   req.AgentVersion,
		Model:          req.Model,
		ModelParams:    req.ModelParams,
		PromptID:       req.PromptID,
		StoreName:      req.StoreName,
		MetadataFilter: req.MetadataFilter,
		Status:         RunStatusPending,
		InputPayload:   req.InputPayload,
		OutputPayload:  nil,
		Error:          nil,
		StartedAt:      nil,
		FinishedAt:     nil,
	}

	return s.repo.CreateGenerationRun(ctx, run)
}

// GetGenerationRunByID retrieves a generation run by ID
func (s *service) GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error) {
	if id == uuid.Nil {
		return nil, errors.New("generation run ID is required")
	}

	return s.repo.GetGenerationRunByID(ctx, id)
}

// ListGenerationRunsByModule retrieves all generation runs for a module
func (s *service) ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListGenerationRunsByModule(ctx, moduleID)
}

// UpdateGenerationRun updates a generation run with business logic validation
func (s *service) UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, error json.RawMessage, startedAt, finishedAt *time.Time) error {
	if id == uuid.Nil {
		return errors.New("generation run ID is required")
	}

	// Check if generation run exists
	_, err := s.repo.GetGenerationRunByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate status transitions if status is being updated
	if status != nil {
		currentRun, err := s.repo.GetGenerationRunByID(ctx, id)
		if err != nil {
			return err
		}

		// Validate status transition rules
		if !isValidStatusTransition(currentRun.Status, *status) {
			return errors.New("invalid status transition")
		}
	}

	return s.repo.UpdateGenerationRun(ctx, id, status, outputPayload, error, startedAt, finishedAt)
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to RunStatus) bool {
	// Define valid status transitions
	validTransitions := map[RunStatus][]RunStatus{
		RunStatusPending:   {RunStatusRunning, RunStatusFailed},
		RunStatusRunning:   {RunStatusCompleted, RunStatusFailed},
		RunStatusCompleted: {},
		RunStatusFailed:    {},
	}

	for _, validTo := range validTransitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// DeleteGenerationRun deletes a generation run
func (s *service) DeleteGenerationRun(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("generation run ID is required")
	}

	// Check if generation run exists before deleting
	_, err := s.repo.GetGenerationRunByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteGenerationRun(ctx, id)
}

// CreateArtifact creates a new artifact with business logic validation
func (s *service) CreateArtifact(ctx context.Context, req CreateArtifactRequest) (*Artifact, error) {
	// Business logic validation
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

	// Create the artifact
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

// GetArtifactByID retrieves an artifact by ID
func (s *service) GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	if id == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	return s.repo.GetArtifactByID(ctx, id)
}

// ListArtifactsByModule retrieves all artifacts for a module
func (s *service) ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListArtifactsByModule(ctx, moduleID)
}

// ListArtifactsByModuleAndStatus retrieves all artifacts for a module with a specific status
func (s *service) ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListArtifactsByModuleAndStatus(ctx, moduleID, status)
}

// UpdateArtifactStatus updates an artifact's status with business logic validation
func (s *service) UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error {
	if id == uuid.Nil {
		return errors.New("artifact ID is required")
	}

	// Check if artifact exists
	currentArtifact, err := s.repo.GetArtifactByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate status transition
	if !isArtifactStatusTransitionValid(currentArtifact.Status, status) {
		return errors.New("invalid artifact status transition")
	}

	// Set timestamps based on status
	now := time.Now()
	if status == ArtifactStatusApproved && approvedAt == nil {
		approvedAt = &now
	} else if status == ArtifactStatusRejected && rejectedAt == nil {
		rejectedAt = &now
	}

	return s.repo.UpdateArtifactStatus(ctx, id, status, approvedAt, rejectedAt)
}

// isArtifactStatusTransitionValid checks if an artifact status transition is valid
func isArtifactStatusTransitionValid(from, to ArtifactStatus) bool {
	// Define valid status transitions for artifacts
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

// DeleteArtifact deletes an artifact
func (s *service) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("artifact ID is required")
	}

	// Check if artifact exists before deleting
	_, err := s.repo.GetArtifactByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteArtifact(ctx, id)
}

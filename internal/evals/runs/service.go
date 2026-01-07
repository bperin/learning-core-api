package runs

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service defines the interface for eval run business logic
type Service interface {
	Create(ctx context.Context, req CreateRunRequest) (*Run, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Run, error)
	GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error)
	ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error)
	UpdateResult(ctx context.Context, id uuid.UUID, req UpdateRunResultRequest) (*Run, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new eval run service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new eval run with business logic validation
func (s *service) Create(ctx context.Context, req CreateRunRequest) (*Run, error) {
	// Business logic validation
	if req.ArtifactID == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	if req.SuiteID == uuid.Nil {
		return nil, errors.New("suite ID is required")
	}

	if req.JudgeModel == "" {
		return nil, errors.New("judge model is required")
	}

	// Create the run
	run := Run{
		ArtifactID:      req.ArtifactID,
		GenerationRunID: req.GenerationRunID,
		SuiteID:         req.SuiteID,
		JudgeModel:      req.JudgeModel,
		JudgeParams:     req.JudgeParams,
		Status:          "RUNNING", // Default status
	}

	return s.repo.Create(ctx, run)
}

// GetByID retrieves an eval run by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Run, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetLatestForArtifact retrieves the latest eval run for an artifact
func (s *service) GetLatestForArtifact(ctx context.Context, artifactID uuid.UUID) (*Run, error) {
	if artifactID == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	return s.repo.GetLatestForArtifact(ctx, artifactID)
}

// ListByArtifact retrieves all eval runs for an artifact
func (s *service) ListByArtifact(ctx context.Context, artifactID uuid.UUID) ([]Run, error) {
	if artifactID == uuid.Nil {
		return nil, errors.New("artifact ID is required")
	}

	return s.repo.ListByArtifact(ctx, artifactID)
}

// UpdateResult updates an eval run's result with business logic validation
func (s *service) UpdateResult(ctx context.Context, id uuid.UUID, req UpdateRunResultRequest) (*Run, error) {
	if id == uuid.Nil {
		return nil, errors.New("eval run ID is required")
	}

	// Check if run exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate status transitions if status is being updated
	if req.Status != "" {
		// Add any status transition validation logic here
	}

	return s.repo.UpdateResult(ctx, id, req.Status, req.OverallPass, req.OverallScore, req.Error)
}

// Delete deletes an eval run
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("eval run ID is required")
	}

	// Check if run exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

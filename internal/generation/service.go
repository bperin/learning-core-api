package generation

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RunService defines business logic for generation runs.
type RunService interface {
	CreateGenerationRun(ctx context.Context, req CreateGenerationRunRequest) (*GenerationRun, error)
	GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error)
	ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error)
	UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, error json.RawMessage, startedAt, finishedAt *time.Time) error
	DeleteGenerationRun(ctx context.Context, id uuid.UUID) error
}

// ArtifactService defines business logic for artifacts.
type ArtifactService interface {
	CreateArtifact(ctx context.Context, req CreateArtifactRequest) (*Artifact, error)
	GetArtifactByID(ctx context.Context, id uuid.UUID) (*Artifact, error)
	ListArtifactsByModule(ctx context.Context, moduleID uuid.UUID) ([]Artifact, error)
	ListArtifactsByModuleAndStatus(ctx context.Context, moduleID uuid.UUID, status ArtifactStatus) ([]Artifact, error)
	UpdateArtifactStatus(ctx context.Context, id uuid.UUID, status ArtifactStatus, approvedAt, rejectedAt *time.Time) error
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
}

// Service bundles generation run and artifact business logic.
type Service interface {
	RunService
	ArtifactService
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

package artifacts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

type Service struct {
	queries *store.Queries
	db      *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		queries: store.New(db),
		db:      db,
	}
}

// CreateArtifactParams matches the structure needed to create an artifact
// using easier-to-use types than the sqlc generated ones (e.g. string vs sql.NullString)
type CreateArtifactParams struct {
	Type             string
	GenerationType   string
	Status           string
	UserID           uuid.UUID
	DocumentID       uuid.NullUUID
	EvalID           uuid.NullUUID
	EvalItemID       uuid.NullUUID
	AttemptID        uuid.NullUUID
	Text             string
	OutputJSON       json.RawMessage
	Model            string
	Prompt           string
	InputHash        string
	PromptTemplateID uuid.NullUUID
	SchemaTemplateID uuid.NullUUID
	PromptRender     string
	ModelParams      json.RawMessage
	Meta             json.RawMessage
	Error            string
}

func (s *Service) CreateArtifact(ctx context.Context, params CreateArtifactParams) (*store.Artifact, error) {
	// Convert standard types to sqlc types
	storeParams := store.CreateArtifactParams{
		Type:             params.Type,
		GenerationType:   toNullGenerationType(params.GenerationType),
		Status:           params.Status,
		EvalID:           params.EvalID,
		EvalItemID:       params.EvalItemID,
		AttemptID:        params.AttemptID,
		ReviewerID:       uuid.NullUUID{}, // Reviewer is usually separate
		Text:             utils.ToNullString(params.Text),
		OutputJson:       utils.ToNullRawMessage(params.OutputJSON),
		Model:            utils.ToNullString(params.Model),
		Prompt:           utils.ToNullString(params.Prompt),
		PromptTemplateID: params.PromptTemplateID,
		SchemaTemplateID: params.SchemaTemplateID,
		ModelParams:      utils.ToNullRawMessage(params.ModelParams),
		PromptRender:     utils.ToNullString(params.PromptRender),
		InputHash:        utils.ToNullString(params.InputHash),
		Meta:             utils.ToNullRawMessage(params.Meta),
		Error:            utils.ToNullString(params.Error),
	}

	artifact, err := s.queries.CreateArtifact(ctx, storeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}

	return &artifact, nil
}

func toNullGenerationType(value string) store.NullGenerationType {
	if value == "" {
		return store.NullGenerationType{}
	}
	return store.NullGenerationType{
		GenerationType: store.GenerationType(value),
		Valid:          true,
	}
}

// ListArtifacts returns paginated artifacts
func (s *Service) ListArtifacts(ctx context.Context, limit, offset int32) ([]store.Artifact, int64, error) {
	artifacts, err := s.queries.ListArtifacts(ctx, store.ListArtifactsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list artifacts: %w", err)
	}

	total, err := s.queries.CountArtifacts(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count artifacts: %w", err)
	}

	return artifacts, total, nil
}

// ListArtifactsByType returns paginated artifacts filtered by type
func (s *Service) ListArtifactsByType(ctx context.Context, artifactType string, limit, offset int32) ([]store.Artifact, int64, error) {
	artifacts, err := s.queries.ListArtifactsByType(ctx, store.ListArtifactsByTypeParams{
		Type:   artifactType,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list artifacts by type: %w", err)
	}

	total, err := s.queries.CountArtifactsByType(ctx, artifactType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count artifacts by type: %w", err)
	}

	return artifacts, total, nil
}

// GetArtifactsByType returns artifacts filtered by type
func (s *Service) GetArtifactsByType(ctx context.Context, artifactType string) ([]store.Artifact, error) {
	artifacts, err := s.queries.GetArtifactsByType(ctx, artifactType)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifacts by type: %w", err)
	}
	return artifacts, nil
}

// GetArtifactsByStatus returns artifacts filtered by status
func (s *Service) GetArtifactsByStatus(ctx context.Context, status string) ([]store.Artifact, error) {
	artifacts, err := s.queries.GetArtifactsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifacts by status: %w", err)
	}
	return artifacts, nil
}

// GetArtifactByID returns a single artifact by ID
func (s *Service) GetArtifactByID(ctx context.Context, id uuid.UUID) (*store.Artifact, error) {
	artifact, err := s.queries.GetArtifact(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}
	return &artifact, nil
}

// GetArtifactStats returns statistics about artifacts
func (s *Service) GetArtifactStats(ctx context.Context) (*store.GetArtifactStatsRow, error) {
	stats, err := s.queries.GetArtifactStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact stats: %w", err)
	}
	return &stats, nil
}

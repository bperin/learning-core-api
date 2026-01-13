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

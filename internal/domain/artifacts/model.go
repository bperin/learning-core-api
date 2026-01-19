package artifacts

import (
	"database/sql"
	"learning-core-api/internal/persistance/store"
	"time"

	"github.com/google/uuid"
)

// Artifact represents an artifact for API responses (Swagger-friendly)
type Artifact struct {
	ID               uuid.UUID  `json:"id"`
	Type             string     `json:"type"`
	Status           string     `json:"status"`
	EvalID           *uuid.UUID `json:"eval_id,omitempty"`
	EvalItemID       *uuid.UUID `json:"eval_item_id,omitempty"`
	AttemptID        *uuid.UUID `json:"attempt_id,omitempty"`
	ReviewerID       *uuid.UUID `json:"reviewer_id,omitempty"`
	Text             *string    `json:"text,omitempty"`
	Model            *string    `json:"model,omitempty"`
	Prompt           *string    `json:"prompt,omitempty"`
	InputHash        *string    `json:"input_hash,omitempty"`
	Error            *string    `json:"error,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	PromptTemplateID *uuid.UUID `json:"prompt_template_id,omitempty"`
}

// ArtifactListResponse represents a paginated response for artifacts
type ArtifactListResponse struct {
	Data       []Artifact     `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination metadata for artifacts
type PaginationMeta struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// ConvertFromStore converts a store.Artifact to domain Artifact
func ConvertFromStore(storeArtifact store.Artifact) Artifact {
	return Artifact{
		ID:               storeArtifact.ID,
		Type:             storeArtifact.Type,
		Status:           storeArtifact.Status,
		EvalID:           toUUIDPtr(storeArtifact.EvalID),
		EvalItemID:       toUUIDPtr(storeArtifact.EvalItemID),
		AttemptID:        toUUIDPtr(storeArtifact.AttemptID),
		ReviewerID:       toUUIDPtr(storeArtifact.ReviewerID),
		Text:             toStringPtr(storeArtifact.Text),
		Model:            toStringPtr(storeArtifact.Model),
		Prompt:           toStringPtr(storeArtifact.Prompt),
		InputHash:        toStringPtr(storeArtifact.InputHash),
		Error:            toStringPtr(storeArtifact.Error),
		CreatedAt:        storeArtifact.CreatedAt,
		PromptTemplateID: toUUIDPtr(storeArtifact.PromptTemplateID),
	}
}

func toUUIDPtr(u uuid.NullUUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}

func toStringPtr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

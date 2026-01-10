package prompt_templates

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines prompt template persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreatePromptTemplateRequest) (*PromptTemplate, error)
	CreateVersion(ctx context.Context, req CreatePromptTemplateVersionRequest) (*PromptTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
	GetActiveByKey(ctx context.Context, key string) (*PromptTemplate, error)
	GetByKeyAndVersion(ctx context.Context, key string, version int32) (*PromptTemplate, error)
	Activate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error)
	DeactivateOtherVersions(ctx context.Context, key string, id uuid.UUID) error
}

package schema_templates

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines schema template persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error)
	CreateVersion(ctx context.Context, req CreateSchemaTemplateVersionRequest) (*SchemaTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
	GetByTypeAndVersion(ctx context.Context, schemaType string, version int32) (*SchemaTemplate, error)
	GetActiveByType(ctx context.Context, schemaType string) (*SchemaTemplate, error)
	ListByType(ctx context.Context, schemaType string) ([]*SchemaTemplate, error)
	ListActive(ctx context.Context) ([]*SchemaTemplate, error)
	ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]*SchemaTemplate, error)
	ListByCurriculum(ctx context.Context, curriculumID uuid.UUID) ([]*SchemaTemplate, error)
	Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
	DeactivateOtherVersions(ctx context.Context, schemaType string, id uuid.UUID) error
}

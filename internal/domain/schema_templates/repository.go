package schema_templates

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines schema template persistence operations.
type Repository interface {
	Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
	GetActiveByType(ctx context.Context, schemaType string) (*SchemaTemplate, error)
	ListByType(ctx context.Context, schemaType string) ([]*SchemaTemplate, error)
	ListActive(ctx context.Context) ([]*SchemaTemplate, error)
	Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error)
}

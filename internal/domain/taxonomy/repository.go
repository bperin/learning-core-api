package taxonomy

import (
	"context"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/documents"
)

// Repository defines taxonomy persistence operations.
type Repository interface {
	CreateNode(ctx context.Context, req CreateTaxonomyNodeRequest) (*TaxonomyNode, error)
	GetByID(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error)
	GetActiveByPath(ctx context.Context, path string) (*TaxonomyNode, error)
	ListByPrefix(ctx context.Context, prefix string) ([]*TaxonomyNode, error)
	Activate(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error)

	CreateDocumentLink(ctx context.Context, req CreateDocumentTaxonomyLinkRequest) (*DocumentTaxonomyLink, error)
	UpdateDocumentLinkState(ctx context.Context, req UpdateDocumentTaxonomyLinkStateRequest) (*DocumentTaxonomyLink, error)

	ListDocumentsByPrefix(ctx context.Context, prefix string) ([]*documents.Document, error)
}

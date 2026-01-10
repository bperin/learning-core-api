package taxonomy

import (
	"context"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/documents"
)

// Service defines taxonomy domain behavior.
type Service interface {
	CreateNode(ctx context.Context, req CreateTaxonomyNodeRequest) (*TaxonomyNode, error)
	GetByID(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error)
	GetActiveByPath(ctx context.Context, path string) (*TaxonomyNode, error)
	ListByPrefix(ctx context.Context, prefix string) ([]*TaxonomyNode, error)
	Activate(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error)

	CreateDocumentLink(ctx context.Context, req CreateDocumentTaxonomyLinkRequest) (*DocumentTaxonomyLink, error)
	UpdateDocumentLinkState(ctx context.Context, req UpdateDocumentTaxonomyLinkStateRequest) (*DocumentTaxonomyLink, error)

	ListDocumentsByPrefix(ctx context.Context, prefix string) ([]*documents.Document, error)
}

type serviceImpl struct {
	repo Repository
}

// NewService creates a new taxonomy service.
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateNode(ctx context.Context, req CreateTaxonomyNodeRequest) (*TaxonomyNode, error) {
	return s.repo.CreateNode(ctx, req)
}

func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *serviceImpl) GetActiveByPath(ctx context.Context, path string) (*TaxonomyNode, error) {
	return s.repo.GetActiveByPath(ctx, path)
}

func (s *serviceImpl) ListByPrefix(ctx context.Context, prefix string) ([]*TaxonomyNode, error) {
	return s.repo.ListByPrefix(ctx, prefix)
}

func (s *serviceImpl) Activate(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error) {
	return s.repo.Activate(ctx, id)
}

func (s *serviceImpl) CreateDocumentLink(ctx context.Context, req CreateDocumentTaxonomyLinkRequest) (*DocumentTaxonomyLink, error) {
	return s.repo.CreateDocumentLink(ctx, req)
}

func (s *serviceImpl) UpdateDocumentLinkState(ctx context.Context, req UpdateDocumentTaxonomyLinkStateRequest) (*DocumentTaxonomyLink, error) {
	return s.repo.UpdateDocumentLinkState(ctx, req)
}

func (s *serviceImpl) ListDocumentsByPrefix(ctx context.Context, prefix string) ([]*documents.Document, error) {
	return s.repo.ListDocumentsByPrefix(ctx, prefix)
}

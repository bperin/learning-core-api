package taxonomy

import (
	"context"

	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

// RepositoryImpl implements Repository using SQLC queries.
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new taxonomy repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// CreateNode creates a taxonomy node.
func (r *RepositoryImpl) CreateNode(ctx context.Context, req CreateTaxonomyNodeRequest) (*TaxonomyNode, error) {
	storeNode, err := r.queries.CreateTaxonomyNode(ctx, store.CreateTaxonomyNodeParams{
		Name:             req.Name,
		Description:      utils.SqlNullString(req.Description),
		ParentID:         utils.PtrToNullUUID(req.ParentID),
		Path:             req.Path,
		Depth:            req.Depth,
		State:            req.State,
		Confidence:       utils.SqlNullFloat64(req.Confidence),
		SourceDocumentID: utils.PtrToNullUUID(req.SourceDocumentID),
		IsActive:         req.IsActive,
		CreatedBy:        utils.PtrToNullUUID(req.CreatedBy),
		ApprovedBy:       utils.PtrToNullUUID(req.ApprovedBy),
		ApprovedAt:       utils.SqlNullTime(req.ApprovedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create taxonomy node: %w", err)
	}

	return toDomainTaxonomyNodeRow(&storeNode), nil
}

// GetByID retrieves a taxonomy node by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error) {
	storeNode, err := r.queries.GetTaxonomyNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get taxonomy node: %w", err)
	}

	return toDomainTaxonomyNode(&storeNode), nil
}

// GetActiveByPath retrieves the active taxonomy node for a path.
func (r *RepositoryImpl) GetActiveByPath(ctx context.Context, path string) (*TaxonomyNode, error) {
	storeNode, err := r.queries.GetActiveTaxonomyNodeByPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get active taxonomy node: %w", err)
	}

	return toDomainTaxonomyNode(&storeNode), nil
}

// ListByPrefix lists taxonomy nodes by prefix.
func (r *RepositoryImpl) ListByPrefix(ctx context.Context, prefix string) ([]*TaxonomyNode, error) {
	storeNodes, err := r.queries.ListTaxonomyNodesByPrefix(ctx, utils.StringToNullString(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to list taxonomy nodes by prefix: %w", err)
	}

	nodes := make([]*TaxonomyNode, len(storeNodes))
	for i, node := range storeNodes {
		nodes[i] = toDomainTaxonomyNode(&node)
	}
	return nodes, nil
}

// Activate marks a taxonomy node as active and deactivates other versions of the same path.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) (*TaxonomyNode, error) {
	storeNode, err := r.queries.ActivateTaxonomyNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate taxonomy node: %w", err)
	}

	return toDomainTaxonomyNodeActivate(&storeNode), nil
}

// CreateDocumentLink creates a document taxonomy link.
func (r *RepositoryImpl) CreateDocumentLink(ctx context.Context, req CreateDocumentTaxonomyLinkRequest) (*DocumentTaxonomyLink, error) {
	storeLink, err := r.queries.CreateDocumentTaxonomyLink(ctx, store.CreateDocumentTaxonomyLinkParams{
		DocumentID:     req.DocumentID,
		TaxonomyNodeID: req.TaxonomyNodeID,
		Confidence:     utils.SqlNullFloat64(req.Confidence),
		State:          req.State,
		ApprovedBy:     utils.PtrToNullUUID(req.ApprovedBy),
		ApprovedAt:     utils.SqlNullTime(req.ApprovedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document taxonomy link: %w", err)
	}

	return toDomainDocumentTaxonomyLink(&storeLink), nil
}

// UpdateDocumentLinkState updates a document taxonomy link state.
func (r *RepositoryImpl) UpdateDocumentLinkState(ctx context.Context, req UpdateDocumentTaxonomyLinkStateRequest) (*DocumentTaxonomyLink, error) {
	storeLink, err := r.queries.UpdateDocumentTaxonomyLinkState(ctx, store.UpdateDocumentTaxonomyLinkStateParams{
		DocumentID:     req.DocumentID,
		TaxonomyNodeID: req.TaxonomyNodeID,
		State:          req.State,
		ApprovedBy:     utils.PtrToNullUUID(req.ApprovedBy),
		ApprovedAt:     utils.SqlNullTime(req.ApprovedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update document taxonomy link: %w", err)
	}

	return toDomainDocumentTaxonomyLink(&storeLink), nil
}

// ListDocumentsByPrefix retrieves documents associated with taxonomy nodes under a prefix.
func (r *RepositoryImpl) ListDocumentsByPrefix(ctx context.Context, prefix string) ([]*documents.Document, error) {
	storeDocs, err := r.queries.ListDocumentsByTaxonomyPrefix(ctx, utils.StringToNullString(prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to list documents by taxonomy prefix: %w", err)
	}

	docs := make([]*documents.Document, len(storeDocs))
	for i, doc := range storeDocs {
		docs[i] = toDomainDocument(doc)
	}
	return docs, nil
}

func toDomainTaxonomyNode(storeNode *store.TaxonomyNode) *TaxonomyNode {
	return &TaxonomyNode{
		ID:               storeNode.ID,
		Name:             storeNode.Name,
		Description:      utils.NullStringToPtr(storeNode.Description),
		ParentID:         utils.NullUUIDToPtr(storeNode.ParentID),
		Path:             storeNode.Path,
		Depth:            storeNode.Depth,
		State:            storeNode.State,
		Confidence:       utils.NullFloat64ToPtr(storeNode.Confidence),
		SourceDocumentID: utils.NullUUIDToPtr(storeNode.SourceDocumentID),
		Version:          storeNode.Version,
		IsActive:         storeNode.IsActive,
		CreatedBy:        utils.NullUUIDToPtr(storeNode.CreatedBy),
		ApprovedBy:       utils.NullUUIDToPtr(storeNode.ApprovedBy),
		ApprovedAt:       utils.NullTimeToPtr(storeNode.ApprovedAt),
		CreatedAt:        storeNode.CreatedAt,
		UpdatedAt:        storeNode.UpdatedAt,
	}
}

func toDomainTaxonomyNodeRow(storeNode *store.CreateTaxonomyNodeRow) *TaxonomyNode {
	return &TaxonomyNode{
		ID:               storeNode.ID,
		Name:             storeNode.Name,
		Description:      utils.NullStringToPtr(storeNode.Description),
		ParentID:         utils.NullUUIDToPtr(storeNode.ParentID),
		Path:             storeNode.Path,
		Depth:            storeNode.Depth,
		State:            storeNode.State,
		Confidence:       utils.NullFloat64ToPtr(storeNode.Confidence),
		SourceDocumentID: utils.NullUUIDToPtr(storeNode.SourceDocumentID),
		Version:          storeNode.Version,
		IsActive:         storeNode.IsActive,
		CreatedBy:        utils.NullUUIDToPtr(storeNode.CreatedBy),
		ApprovedBy:       utils.NullUUIDToPtr(storeNode.ApprovedBy),
		ApprovedAt:       utils.NullTimeToPtr(storeNode.ApprovedAt),
		CreatedAt:        storeNode.CreatedAt,
		UpdatedAt:        storeNode.UpdatedAt,
	}
}

func toDomainTaxonomyNodeActivate(storeNode *store.ActivateTaxonomyNodeRow) *TaxonomyNode {
	return &TaxonomyNode{
		ID:               storeNode.ID,
		Name:             storeNode.Name,
		Description:      utils.NullStringToPtr(storeNode.Description),
		ParentID:         utils.NullUUIDToPtr(storeNode.ParentID),
		Path:             storeNode.Path,
		Depth:            storeNode.Depth,
		State:            storeNode.State,
		Confidence:       utils.NullFloat64ToPtr(storeNode.Confidence),
		SourceDocumentID: utils.NullUUIDToPtr(storeNode.SourceDocumentID),
		Version:          storeNode.Version,
		IsActive:         storeNode.IsActive,
		CreatedBy:        utils.NullUUIDToPtr(storeNode.CreatedBy),
		ApprovedBy:       utils.NullUUIDToPtr(storeNode.ApprovedBy),
		ApprovedAt:       utils.NullTimeToPtr(storeNode.ApprovedAt),
		CreatedAt:        storeNode.CreatedAt,
		UpdatedAt:        storeNode.UpdatedAt,
	}
}

func toDomainDocumentTaxonomyLink(storeLink *store.DocumentTaxonomyLink) *DocumentTaxonomyLink {
	return &DocumentTaxonomyLink{
		DocumentID:     storeLink.DocumentID,
		TaxonomyNodeID: storeLink.TaxonomyNodeID,
		Confidence:     utils.NullFloat64ToPtr(storeLink.Confidence),
		State:          storeLink.State,
		CreatedAt:      storeLink.CreatedAt,
		ApprovedBy:     utils.NullUUIDToPtr(storeLink.ApprovedBy),
		ApprovedAt:     utils.NullTimeToPtr(storeLink.ApprovedAt),
	}
}

func toDomainDocument(storeDoc store.Document) *documents.Document {
	return &documents.Document{
		ID:                storeDoc.ID,
		Filename:          storeDoc.Filename,
		Title:             utils.NullStringToPtr(storeDoc.Title),
		MimeType:          utils.NullStringToPtr(storeDoc.MimeType),
		Content:           utils.NullStringToPtr(storeDoc.Content),
		StoragePath:       utils.NullStringToPtr(storeDoc.StoragePath),
		StorageBucket:     utils.NullStringToPtr(storeDoc.StorageBucket),
		FileStoreName:     utils.NullStringToPtr(storeDoc.FileStoreName),
		FileStoreFileName: utils.NullStringToPtr(storeDoc.FileStoreFileName),
		RagStatus:         storeDoc.RagStatus,
		UserID:            storeDoc.UserID,
		CreatedAt:         storeDoc.CreatedAt,
		UpdatedAt:         storeDoc.UpdatedAt,
	}
}

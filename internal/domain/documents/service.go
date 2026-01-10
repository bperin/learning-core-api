package documents

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service defines the business logic interface for documents
type Service interface {
	// CreateDocument creates a new document with business logic validation
	CreateDocument(ctx context.Context, req CreateDocumentRequest) (*Document, error)

	// GetDocument retrieves a document by ID with authorization checks
	GetDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Document, error)

	// GetUserDocuments retrieves all documents for a user
	GetUserDocuments(ctx context.Context, userID uuid.UUID) ([]*Document, error)

	// SearchDocuments searches documents with authorization
	SearchDocuments(ctx context.Context, query string, userID uuid.UUID, limit, offset int) ([]*Document, error)

	// UpdateDocument updates a document with authorization and validation
	UpdateDocument(ctx context.Context, id uuid.UUID, req UpdateDocumentRequest, userID uuid.UUID) (*Document, error)

	// UpdateDocumentRagStatus updates the RAG processing status
	UpdateDocumentRagStatus(ctx context.Context, id uuid.UUID, status string) (*Document, error)

	// DeleteDocument deletes a document with authorization
	DeleteDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// ProcessDocument initiates document processing for RAG
	ProcessDocument(ctx context.Context, id uuid.UUID) error

	// GetDocumentsByStatus retrieves documents by RAG status (admin only)
	GetDocumentsByStatus(ctx context.Context, status string) ([]*Document, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	repo Repository
}

// NewService creates a new document service
func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

// CreateDocument creates a new document with business logic validation
func (s *serviceImpl) CreateDocument(ctx context.Context, req CreateDocumentRequest) (*Document, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Additional business logic validation
	if err := s.validateFileType(req.MimeType); err != nil {
		return nil, err
	}

	// Set default RAG status if not provided
	if req.RagStatus == "" {
		req.RagStatus = RagStatusPending
	}

	// Create the document
	doc, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return doc, nil
}

// GetDocument retrieves a document by ID with authorization checks
func (s *serviceImpl) GetDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Document, error) {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Authorization check - user can only access their own documents
	if doc.UserID != userID {
		return nil, ErrUnauthorized
	}

	return doc, nil
}

// GetUserDocuments retrieves all documents for a user
func (s *serviceImpl) GetUserDocuments(ctx context.Context, userID uuid.UUID) ([]*Document, error) {
	return s.repo.GetByUser(ctx, userID)
}

// SearchDocuments searches documents with authorization
func (s *serviceImpl) SearchDocuments(ctx context.Context, query string, userID uuid.UUID, limit, offset int) ([]*Document, error) {
	docs, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	// Filter documents to only include those owned by the user
	var userDocs []*Document
	for _, doc := range docs {
		if doc.UserID == userID {
			userDocs = append(userDocs, doc)
		}
	}

	return userDocs, nil
}

// UpdateDocument updates a document with authorization and validation
func (s *serviceImpl) UpdateDocument(ctx context.Context, id uuid.UUID, req UpdateDocumentRequest, userID uuid.UUID) (*Document, error) {
	// First check if the document exists and user has access
	existingDoc, err := s.GetDocument(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Validate RAG status if being updated
	if req.RagStatus != nil && !IsValidRagStatus(*req.RagStatus) {
		return nil, ErrInvalidRagStatus
	}

	// Additional business logic validation
	// Note: MimeType is not updateable in UpdateDocumentRequest

	// Prevent updating certain fields if document is being processed
	if existingDoc.RagStatus == RagStatusProcessing {
		if req.Content != nil || req.StoragePath != nil {
			return nil, ErrProcessingFailed
		}
	}

	return s.repo.Update(ctx, id, req)
}

// UpdateDocumentRagStatus updates the RAG processing status
func (s *serviceImpl) UpdateDocumentRagStatus(ctx context.Context, id uuid.UUID, status string) (*Document, error) {
	if !IsValidRagStatus(status) {
		return nil, ErrInvalidRagStatus
	}

	return s.repo.UpdateRagStatus(ctx, id, status)
}

// DeleteDocument deletes a document with authorization
func (s *serviceImpl) DeleteDocument(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// First check if the document exists and user has access
	_, err := s.GetDocument(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// ProcessDocument initiates document processing for RAG
func (s *serviceImpl) ProcessDocument(ctx context.Context, id uuid.UUID) error {
	// Update status to processing
	_, err := s.repo.UpdateRagStatus(ctx, id, RagStatusProcessing)
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// TODO: Trigger actual RAG processing pipeline
	// This would typically involve:
	// 1. Extracting text from the document
	// 2. Chunking the text
	// 3. Generating embeddings
	// 4. Storing in vector database
	// 5. Updating status to ready or error

	return nil
}

// GetDocumentsByStatus retrieves documents by RAG status (admin only)
func (s *serviceImpl) GetDocumentsByStatus(ctx context.Context, status string) ([]*Document, error) {
	if !IsValidRagStatus(status) {
		return nil, ErrInvalidRagStatus
	}

	return s.repo.GetByRagStatus(ctx, status)
}

// validateFileType validates the MIME type of the document
func (s *serviceImpl) validateFileType(mimeType *string) error {
	if mimeType == nil {
		return nil // Optional field
	}

	allowedTypes := map[string]bool{
		"application/pdf": true,
		"text/plain":      true,
		"text/markdown":   true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
		"application/msword": true, // .doc
		"text/html":          true,
		"application/json":   true,
	}

	if !allowedTypes[*mimeType] {
		return ErrInvalidFileType
	}

	return nil
}

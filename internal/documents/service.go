package documents

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
)

// Service defines the interface for document business logic
type Service interface {
	Create(ctx context.Context, req CreateDocumentRequest) (*Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Document, error)
	GetByModuleAndSourceURI(ctx context.Context, moduleID uuid.UUID, sourceURI string) (*Document, error)
	ListByModule(ctx context.Context, moduleID uuid.UUID) ([]Document, error)
	Update(ctx context.Context, id uuid.UUID, title *string, metadata map[string]interface{}, indexedAt *time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new document service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Create creates a new document with business logic validation
func (s *service) Create(ctx context.Context, req CreateDocumentRequest) (*Document, error) {
	// Business logic validation
	if req.ModuleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if req.SourceURI == "" {
		return nil, errors.New("source URI is required")
	}

	// Check if document with same module and source URI already exists
	existing, err := s.repo.GetByModuleAndSourceURI(ctx, req.ModuleID, req.SourceURI)
	if err == nil && existing != nil {
		return nil, errors.New("document with this source URI already exists for the module")
	}

	// Create the document
	document := Document{
		ModuleID:   req.ModuleID,
		StoreID:    req.StoreID,
		Title:      req.Title,
		SourceURI:  req.SourceURI,
		SHA256:     req.SHA256,
		Metadata:   req.Metadata,
		FileName:   req.FileName,
		DocName:    req.DocName,
		IndexedAt:  req.IndexedAt,
	}

	return s.repo.Create(ctx, document)
}

// GetByID retrieves a document by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Document, error) {
	if id == uuid.Nil {
		return nil, errors.New("document ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetByModuleAndSourceURI retrieves a document by module ID and source URI
func (s *service) GetByModuleAndSourceURI(ctx context.Context, moduleID uuid.UUID, sourceURI string) (*Document, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if sourceURI == "" {
		return nil, errors.New("source URI is required")
	}

	return s.repo.GetByModuleAndSourceURI(ctx, moduleID, sourceURI)
}

// ListByModule retrieves all documents for a module
func (s *service) ListByModule(ctx context.Context, moduleID uuid.UUID) ([]Document, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListByModule(ctx, moduleID)
}

// Update updates a document with business logic validation
func (s *service) Update(ctx context.Context, id uuid.UUID, title *string, metadata map[string]interface{}, indexedAt *time.Time) error {
	if id == uuid.Nil {
		return errors.New("document ID is required")
	}

	// Check if document exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, id, title, metadata, indexedAt)
}

// Delete deletes a document
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("document ID is required")
	}

	// Check if document exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
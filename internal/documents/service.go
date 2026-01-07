package documents

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for document business logic
type Service interface {
	Create(ctx context.Context, req CreateDocumentRequest) (*Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Document, error)
	GetBySubjectAndSourceURI(ctx context.Context, subjectID uuid.UUID, sourceURI string) (*Document, error)
	ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]Document, error)
	Update(ctx context.Context, id uuid.UUID, title *string, metadata json.RawMessage, indexedAt *time.Time) error
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
	if req.SubjectID == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}

	if req.SourceURI == "" {
		return nil, errors.New("source URI is required")
	}

	// Check if document with same module and source URI already exists
	existing, err := s.repo.GetBySubjectAndSourceURI(ctx, req.SubjectID, req.SourceURI)
	if err == nil && existing != nil {
		return nil, errors.New("document with this source URI already exists for the subject")
	}

	// Create the document
	document := Document{
		SubjectID: req.SubjectID,
		StoreID:   req.StoreID,
		Title:     req.Title,
		SourceURI: req.SourceURI,
		SHA256:    req.SHA256,
		Metadata:  req.Metadata,
		FileName:  req.FileName,
		DocName:   req.DocName,
		IndexedAt: req.IndexedAt,
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

// GetBySubjectAndSourceURI retrieves a document by subject ID and source URI.
func (s *service) GetBySubjectAndSourceURI(ctx context.Context, subjectID uuid.UUID, sourceURI string) (*Document, error) {
	if subjectID == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}

	if sourceURI == "" {
		return nil, errors.New("source URI is required")
	}

	return s.repo.GetBySubjectAndSourceURI(ctx, subjectID, sourceURI)
}

// ListBySubject retrieves all documents for a subject.
func (s *service) ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]Document, error) {
	if subjectID == uuid.Nil {
		return nil, errors.New("subject ID is required")
	}

	return s.repo.ListBySubject(ctx, subjectID)
}

// Update updates a document with business logic validation
func (s *service) Update(ctx context.Context, id uuid.UUID, title *string, metadata json.RawMessage, indexedAt *time.Time) error {
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

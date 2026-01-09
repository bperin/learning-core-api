package documents

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for document data access
type Repository interface {
	// Create creates a new document
	Create(ctx context.Context, req CreateDocumentRequest) (*Document, error)

	// GetByID retrieves a document by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Document, error)

	// GetByUser retrieves all documents for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*Document, error)

	// GetBySubject retrieves all documents in a specific subject
	GetBySubject(ctx context.Context, subjectID uuid.UUID) ([]*Document, error)

	// GetByRagStatus retrieves documents by RAG processing status
	GetByRagStatus(ctx context.Context, status string) ([]*Document, error)

	// GetBySubjects retrieves documents that match any of the provided subjects
	GetBySubjects(ctx context.Context, subjects []string) ([]*Document, error)

	// List retrieves documents with pagination
	List(ctx context.Context, limit, offset int) ([]*Document, error)

	// Search searches documents by title with pagination
	Search(ctx context.Context, title string, limit, offset int) ([]*Document, error)

	// Update updates an existing document
	Update(ctx context.Context, id uuid.UUID, req UpdateDocumentRequest) (*Document, error)

	// UpdateRagStatus updates only the RAG status of a document
	UpdateRagStatus(ctx context.Context, id uuid.UUID, status string) (*Document, error)

	// Delete deletes a document by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// Filter retrieves documents based on filter criteria
	Filter(ctx context.Context, filter DocumentFilter) ([]*Document, error)
}

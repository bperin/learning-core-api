package documents

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

// repositoryImpl implements the Repository interface using sqlc generated code
type repositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new document repository
func NewRepository(queries *store.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

// Create creates a new document
func (r *repositoryImpl) Create(ctx context.Context, req CreateDocumentRequest) (*Document, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Set defaults
	if req.RagStatus == "" {
		req.RagStatus = RagStatusPending
	}
	if req.Subjects == nil {
		req.Subjects = []string{}
	}

	params := store.CreateDocumentParams{
		Filename:    req.Filename,
		Title:       utils.SqlNullString(req.Title),
		MimeType:    utils.SqlNullString(req.MimeType),
		Content:     utils.SqlNullString(req.Content),
		StoragePath: utils.SqlNullString(req.StoragePath),
		RagStatus:   req.RagStatus,
		UserID:      req.UserID,
		SubjectID:   utils.SqlNullUUID(req.SubjectID),
		Curricular:  utils.SqlNullString(req.Curricular),
		Subjects:    req.Subjects,
	}

	doc, err := r.queries.CreateDocument(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toDomain(doc), nil
}

// GetByID retrieves a document by its ID
func (r *repositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*Document, error) {
	doc, err := r.queries.GetDocument(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	return r.toDomain(doc), nil
}

// GetByUser retrieves all documents for a specific user
func (r *repositoryImpl) GetByUser(ctx context.Context, userID uuid.UUID) ([]*Document, error) {
	docs, err := r.queries.GetDocumentsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// GetBySubject retrieves all documents in a specific subject
func (r *repositoryImpl) GetBySubject(ctx context.Context, subjectID uuid.UUID) ([]*Document, error) {
	docs, err := r.queries.GetDocumentsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// GetByRagStatus retrieves documents by RAG processing status
func (r *repositoryImpl) GetByRagStatus(ctx context.Context, status string) ([]*Document, error) {
	docs, err := r.queries.GetDocumentsByRagStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// GetBySubjects retrieves documents that match any of the provided subjects
func (r *repositoryImpl) GetBySubjects(ctx context.Context, subjects []string) ([]*Document, error) {
	docs, err := r.queries.GetDocumentsBySubjects(ctx, subjects)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// List retrieves documents with pagination
func (r *repositoryImpl) List(ctx context.Context, limit, offset int) ([]*Document, error) {
	params := store.ListDocumentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	docs, err := r.queries.ListDocuments(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// Search searches documents by title with pagination
func (r *repositoryImpl) Search(ctx context.Context, title string, limit, offset int) ([]*Document, error) {
	params := store.SearchDocumentsByTitleParams{
		Title:      utils.SqlNullString(&title),
		PageLimit:  int32(limit),
		PageOffset: int32(offset),
	}

	docs, err := r.queries.SearchDocumentsByTitle(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(docs), nil
}

// Update updates an existing document
func (r *repositoryImpl) Update(ctx context.Context, id uuid.UUID, req UpdateDocumentRequest) (*Document, error) {
	params := store.UpdateDocumentParams{
		ID:          id,
		Title:       utils.SqlNullString(req.Title),
		Content:     utils.SqlNullString(req.Content),
		StoragePath: utils.SqlNullString(req.StoragePath),
		RagStatus:   utils.NullStringToString(utils.SqlNullString(req.RagStatus)),
		SubjectID:   utils.SqlNullUUID(req.SubjectID),
		Curricular:  utils.SqlNullString(req.Curricular),
		Subjects:    req.Subjects,
	}

	doc, err := r.queries.UpdateDocument(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	return r.toDomain(doc), nil
}

// UpdateRagStatus updates only the RAG status of a document
func (r *repositoryImpl) UpdateRagStatus(ctx context.Context, id uuid.UUID, status string) (*Document, error) {
	if !IsValidRagStatus(status) {
		return nil, ErrInvalidRagStatus
	}

	params := store.UpdateDocumentRagStatusParams{
		ID:        id,
		RagStatus: status,
	}

	doc, err := r.queries.UpdateDocumentRagStatus(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	return r.toDomain(doc), nil
}

// Delete deletes a document by ID
func (r *repositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteDocument(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// Filter retrieves documents based on filter criteria
func (r *repositoryImpl) Filter(ctx context.Context, filter DocumentFilter) ([]*Document, error) {
	// For complex filtering, we'll use multiple queries based on the filter criteria
	// This is a simplified implementation - in production you might want a more sophisticated approach

	if filter.UserID != nil {
		return r.GetByUser(ctx, *filter.UserID)
	}

	if filter.SubjectID != nil {
		return r.GetBySubject(ctx, *filter.SubjectID)
	}

	if filter.RagStatus != nil {
		return r.GetByRagStatus(ctx, *filter.RagStatus)
	}

	if len(filter.Subjects) > 0 {
		return r.GetBySubjects(ctx, filter.Subjects)
	}

	if filter.Title != nil {
		return r.Search(ctx, *filter.Title, filter.Limit, filter.Offset)
	}

	// Default to list with pagination
	return r.List(ctx, filter.Limit, filter.Offset)
}

// Helper functions for converting between domain and store types

func (r *repositoryImpl) toDomain(doc store.Document) *Document {
	return &Document{
		ID:          doc.ID,
		Filename:    doc.Filename,
		Title:       utils.NullStringToPtr(doc.Title),
		MimeType:    utils.NullStringToPtr(doc.MimeType),
		Content:     utils.NullStringToPtr(doc.Content),
		StoragePath: utils.NullStringToPtr(doc.StoragePath),
		RagStatus:   doc.RagStatus,
		UserID:      doc.UserID,
		SubjectID:   utils.NullUUIDToPtr(doc.SubjectID),
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
		Curricular:  utils.NullStringToPtr(doc.Curricular),
		Subjects:    doc.Subjects,
	}
}

func (r *repositoryImpl) toDomainSlice(docs []store.Document) []*Document {
	result := make([]*Document, len(docs))
	for i, doc := range docs {
		result[i] = r.toDomain(doc)
	}
	return result
}

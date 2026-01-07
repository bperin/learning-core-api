package documents

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"learning-core-api/internal/store"
	"learning-core-api/internal/utils"
)

// Repository defines the interface for document operations
type Repository interface {
	Create(ctx context.Context, document Document) (*Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Document, error)
	GetByModuleAndSourceURI(ctx context.Context, moduleID uuid.UUID, sourceURI string) (*Document, error)
	ListByModule(ctx context.Context, moduleID uuid.UUID) ([]Document, error)
	Update(ctx context.Context, id uuid.UUID, title *string, metadata map[string]interface{}, indexedAt *time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new document repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Helper functions
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func bytesToMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	if len(data) > 0 {
		json.Unmarshal(data, &result)
	}
	return result
}

// Create creates a new document
func (r *repository) Create(ctx context.Context, document Document) (*Document, error) {
	params := store.CreateDocumentParams{
		ModuleID:  document.ModuleID,
		StoreID:   document.StoreID,
		SourceUri: document.SourceURI,
	}

	// Handle nullable fields
	if document.Title != "" {
		params.Title = sql.NullString{String: document.Title, Valid: true}
	} else {
		params.Title = sql.NullString{String: "", Valid: false}
	}

	if document.SHA256 != "" {
		params.Sha256 = sql.NullString{String: document.SHA256, Valid: true}
	} else {
		params.Sha256 = sql.NullString{String: "", Valid: false}
	}

	if document.FileName != "" {
		params.FileName = sql.NullString{String: document.FileName, Valid: true}
	} else {
		params.FileName = sql.NullString{String: "", Valid: false}
	}

	if document.DocName != "" {
		params.DocName = sql.NullString{String: document.DocName, Valid: true}
	} else {
		params.DocName = sql.NullString{String: "", Valid: false}
	}

	// Handle JSON metadata
	metadataBytes, err := json.Marshal(document.Metadata)
	if err != nil {
		return nil, err
	}
	params.Metadata = metadataBytes

	dbDocument, err := r.queries.CreateDocument(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Document{
		ID:        dbDocument.ID,
		ModuleID:  dbDocument.ModuleID,
		StoreID:   dbDocument.StoreID,
		Title:     nullStringToString(dbDocument.Title),
		SourceURI: dbDocument.SourceUri,
		SHA256:    nullStringToString(dbDocument.Sha256),
		Metadata:  bytesToMap(dbDocument.Metadata),
		FileName:  nullStringToString(dbDocument.FileName),
		DocName:   nullStringToString(dbDocument.DocName),
		IndexedAt: utils.NullTimeToPtr(dbDocument.IndexedAt),
		CreatedAt: dbDocument.CreatedAt,
	}, nil
}

// GetByID retrieves a document by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Document, error) {
	dbDocument, err := r.queries.GetDocument(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Document{
		ID:        dbDocument.ID,
		ModuleID:  dbDocument.ModuleID,
		StoreID:   dbDocument.StoreID,
		Title:     nullStringToString(dbDocument.Title),
		SourceURI: dbDocument.SourceUri,
		SHA256:    nullStringToString(dbDocument.Sha256),
		Metadata:  bytesToMap(dbDocument.Metadata),
		FileName:  nullStringToString(dbDocument.FileName),
		DocName:   nullStringToString(dbDocument.DocName),
		IndexedAt: utils.NullTimeToPtr(dbDocument.IndexedAt),
		CreatedAt: dbDocument.CreatedAt,
	}, nil
}

// GetByModuleAndSourceURI retrieves a document by module ID and source URI
func (r *repository) GetByModuleAndSourceURI(ctx context.Context, moduleID uuid.UUID, sourceURI string) (*Document, error) {
	dbDocument, err := r.queries.GetDocumentByModuleAndSourceURI(ctx, store.GetDocumentByModuleAndSourceURIParams{
		ModuleID:  moduleID,
		SourceUri: sourceURI,
	})
	if err != nil {
		return nil, err
	}

	return &Document{
		ID:        dbDocument.ID,
		ModuleID:  dbDocument.ModuleID,
		StoreID:   dbDocument.StoreID,
		Title:     nullStringToString(dbDocument.Title),
		SourceURI: dbDocument.SourceUri,
		SHA256:    nullStringToString(dbDocument.Sha256),
		Metadata:  bytesToMap(dbDocument.Metadata),
		FileName:  nullStringToString(dbDocument.FileName),
		DocName:   nullStringToString(dbDocument.DocName),
		IndexedAt: utils.NullTimeToPtr(dbDocument.IndexedAt),
		CreatedAt: dbDocument.CreatedAt,
	}, nil
}

// ListByModule retrieves all documents for a module
func (r *repository) ListByModule(ctx context.Context, moduleID uuid.UUID) ([]Document, error) {
	dbDocuments, err := r.queries.ListDocumentsByModule(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	documents := make([]Document, len(dbDocuments))
	for i, dbDocument := range dbDocuments {
		documents[i] = Document{
			ID:        dbDocument.ID,
			ModuleID:  dbDocument.ModuleID,
			StoreID:   dbDocument.StoreID,
			Title:     nullStringToString(dbDocument.Title),
			SourceURI: dbDocument.SourceUri,
			SHA256:    nullStringToString(dbDocument.Sha256),
			Metadata:  bytesToMap(dbDocument.Metadata),
			FileName:  nullStringToString(dbDocument.FileName),
			DocName:   nullStringToString(dbDocument.DocName),
			IndexedAt: utils.NullTimeToPtr(dbDocument.IndexedAt),
			CreatedAt: dbDocument.CreatedAt,
		}
	}

	return documents, nil
}

// Update updates a document
func (r *repository) Update(ctx context.Context, id uuid.UUID, title *string, metadata map[string]interface{}, indexedAt *time.Time) error {
	// Get current document to preserve values if not updating
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	params := store.UpdateDocumentParams{
		ID: id,
	}

	// Set title parameter
	if title != nil {
		params.Title = sql.NullString{String: *title, Valid: true}
	} else {
		params.Title = sql.NullString{String: current.Title, Valid: current.Title != ""}
	}

	// Set metadata parameter
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		params.Metadata = metadataBytes
	} else {
		metadataBytes, err := json.Marshal(current.Metadata)
		if err != nil {
			return err
		}
		params.Metadata = metadataBytes
	}

	// Set indexedAt parameter
	if indexedAt != nil {
		params.IndexedAt = sql.NullTime{Time: *indexedAt, Valid: true}
	} else if current.IndexedAt != nil {
		params.IndexedAt = sql.NullTime{Time: *current.IndexedAt, Valid: true}
	} else {
		params.IndexedAt = sql.NullTime{Valid: false}
	}

	return r.queries.UpdateDocument(ctx, params)
}

// Delete deletes a document
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteDocument(ctx, id)
}

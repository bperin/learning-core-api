package documents

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Document represents a document in the learning platform
type Document struct {
	ID        uuid.UUID       `json:"id"`
	SubjectID uuid.UUID       `json:"subject_id"`
	StoreID   uuid.UUID       `json:"store_id"`
	Title     string          `json:"title"`
	SourceURI string          `json:"source_uri"`
	SHA256    string          `json:"sha256"`
	Metadata  json.RawMessage `json:"metadata"`
	FileName  string          `json:"file_name"`
	DocName   string          `json:"doc_name"`
	IndexedAt *time.Time      `json:"indexed_at,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// CreateDocumentRequest represents the request to create a document
type CreateDocumentRequest struct {
	SubjectID uuid.UUID       `json:"subject_id"`
	StoreID   uuid.UUID       `json:"store_id"`
	Title     string          `json:"title"`
	SourceURI string          `json:"source_uri"`
	SHA256    string          `json:"sha256"`
	Metadata  json.RawMessage `json:"metadata"`
	FileName  string          `json:"file_name"`
	DocName   string          `json:"doc_name"`
	IndexedAt *time.Time      `json:"indexed_at,omitempty"`
}

// UpdateDocumentRequest represents the request to update a document
type UpdateDocumentRequest struct {
	Title     *string         `json:"title,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	IndexedAt *time.Time      `json:"indexed_at,omitempty"`
}

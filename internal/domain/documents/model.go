package documents

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a learning document in the domain
type Document struct {
	ID                uuid.UUID `json:"id"`
	Filename          string    `json:"filename"`
	Title             *string   `json:"title,omitempty"`
	MimeType          *string   `json:"mime_type,omitempty"`
	Content           *string   `json:"content,omitempty"`
	StoragePath       *string   `json:"storage_path,omitempty"`
	StorageBucket     *string   `json:"storage_bucket,omitempty"`
	FileStoreName     *string   `json:"file_store_name,omitempty"`
	FileStoreFileName *string   `json:"file_store_file_name,omitempty"`
	RagStatus         string    `json:"rag_status"`
	UserID            uuid.UUID `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateDocumentRequest represents the data needed to create a new document
type CreateDocumentRequest struct {
	Filename          string    `json:"filename" validate:"required"`
	Title             *string   `json:"title,omitempty"`
	MimeType          *string   `json:"mime_type,omitempty"`
	Content           *string   `json:"content,omitempty"`
	StoragePath       *string   `json:"storage_path,omitempty"`
	StorageBucket     *string   `json:"storage_bucket,omitempty"`
	FileStoreName     *string   `json:"file_store_name,omitempty"`
	FileStoreFileName *string   `json:"file_store_file_name,omitempty"`
	RagStatus         string    `json:"rag_status" validate:"required"`
	UserID            uuid.UUID `json:"user_id" validate:"required"`
}

// UpdateDocumentRequest represents the data that can be updated for a document
type UpdateDocumentRequest struct {
	Title             *string `json:"title,omitempty"`
	Content           *string `json:"content,omitempty"`
	StoragePath       *string `json:"storage_path,omitempty"`
	StorageBucket     *string `json:"storage_bucket,omitempty"`
	FileStoreName     *string `json:"file_store_name,omitempty"`
	FileStoreFileName *string `json:"file_store_file_name,omitempty"`
	RagStatus         *string `json:"rag_status,omitempty"`
}

// DocumentFilter represents filtering options for document queries
type DocumentFilter struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	RagStatus *string    `json:"rag_status,omitempty"`
	Title     *string    `json:"title,omitempty"` // For search
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// RAG Status constants
const (
	RagStatusPending    = "pending"
	RagStatusProcessing = "processing"
	RagStatusReady      = "ready"
	RagStatusError      = "error"
)

// IsValidRagStatus checks if the provided status is valid
func IsValidRagStatus(status string) bool {
	switch status {
	case RagStatusPending, RagStatusProcessing, RagStatusReady, RagStatusError:
		return true
	default:
		return false
	}
}

// SetDefaults sets default values for the document
func (d *Document) SetDefaults() {
	if d.RagStatus == "" {
		d.RagStatus = RagStatusPending
	}
}

// Validate performs domain-level validation
func (req *CreateDocumentRequest) Validate() error {
	if req.Filename == "" {
		return ErrInvalidFilename
	}
	if req.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if req.RagStatus != "" && !IsValidRagStatus(req.RagStatus) {
		return ErrInvalidRagStatus
	}
	return nil
}

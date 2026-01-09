package documents

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a learning document in the domain
type Document struct {
	ID          uuid.UUID  `json:"id"`
	Filename    string     `json:"filename"`
	Title       *string    `json:"title,omitempty"`
	MimeType    *string    `json:"mime_type,omitempty"`
	Content     *string    `json:"content,omitempty"`
	StoragePath *string    `json:"storage_path,omitempty"`
	RagStatus   string     `json:"rag_status"`
	UserID      uuid.UUID  `json:"user_id"`
	SubjectID   *uuid.UUID `json:"subject_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	// Curricular classification or framework (e.g., Common Core, IB)
	Curricular *string `json:"curricular,omitempty"`
	// List of academic subjects associated with this document
	Subjects []string `json:"subjects"`
}

// CreateDocumentRequest represents the data needed to create a new document
type CreateDocumentRequest struct {
	Filename    string     `json:"filename" validate:"required"`
	Title       *string    `json:"title,omitempty"`
	MimeType    *string    `json:"mime_type,omitempty"`
	Content     *string    `json:"content,omitempty"`
	StoragePath *string    `json:"storage_path,omitempty"`
	RagStatus   string     `json:"rag_status" validate:"required"`
	UserID      uuid.UUID  `json:"user_id" validate:"required"`
	SubjectID   *uuid.UUID `json:"subject_id,omitempty"`
	Curricular  *string    `json:"curricular,omitempty"`
	Subjects    []string   `json:"subjects"`
}

// UpdateDocumentRequest represents the data that can be updated for a document
type UpdateDocumentRequest struct {
	Title       *string    `json:"title,omitempty"`
	Content     *string    `json:"content,omitempty"`
	StoragePath *string    `json:"storage_path,omitempty"`
	RagStatus   *string    `json:"rag_status,omitempty"`
	SubjectID   *uuid.UUID `json:"subject_id,omitempty"`
	Curricular  *string    `json:"curricular,omitempty"`
	Subjects    []string   `json:"subjects,omitempty"`
}

// DocumentFilter represents filtering options for document queries
type DocumentFilter struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	SubjectID *uuid.UUID `json:"subject_id,omitempty"`
	RagStatus *string    `json:"rag_status,omitempty"`
	Subjects  []string   `json:"subjects,omitempty"`
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
	if d.Subjects == nil {
		d.Subjects = []string{}
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

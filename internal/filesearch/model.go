package filesearch

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Store represents a file search store tied to a subject.
type Store struct {
	ID             uuid.UUID       `json:"id"`
	SubjectID      uuid.UUID       `json:"subject_id"`
	StoreName      string          `json:"store_name"`
	DisplayName    string          `json:"display_name"`
	ChunkingConfig json.RawMessage `json:"chunking_config"`
	CreatedAt      time.Time       `json:"created_at"`
}

// UploadRequest captures parameters for uploading a file to a subject store.
type UploadRequest struct {
	SubjectID   uuid.UUID       `json:"subject_id"`
	FilePath    string          `json:"file_path"`
	SourceURI   string          `json:"source_uri"`
	Title       string          `json:"title"`
	FileName    string          `json:"file_name"`
	DocName     string          `json:"doc_name"`
	DisplayName string          `json:"display_name"`
	MimeType    string          `json:"mime_type"`
	Metadata    json.RawMessage `json:"metadata"`
}

// UploadResult captures upload and document creation results.
type UploadResult struct {
	Store        *Store `json:"store"`
	DocumentName string `json:"document_name"`
}

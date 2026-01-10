package documents

import "errors"

// Domain errors for documents
var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidFilename  = errors.New("invalid filename")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidRagStatus = errors.New("invalid RAG status")
	ErrDocumentExists   = errors.New("document already exists")
	ErrUnauthorized     = errors.New("unauthorized access to document")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrFileTooLarge     = errors.New("file too large")
	ErrProcessingFailed = errors.New("document processing failed")
	ErrStorageError     = errors.New("storage error")
)

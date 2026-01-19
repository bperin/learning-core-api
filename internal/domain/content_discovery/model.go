package content_discovery

import (
	"github.com/google/uuid"
)

// Book represents a book found from a subject URL
type Book struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	PDFLink     string `json:"pdf_link,omitempty"`
	Description string `json:"description,omitempty"`
	Authors     string `json:"authors,omitempty"`
	Selected    bool   `json:"selected"`
}

// BookListRequest represents a request to list books from subjects
type BookListRequest struct {
	SubjectIDs []string `json:"subject_ids"`
	MaxBooks   int      `json:"max_books,omitempty"` // Default to 10 if not specified
}

// BookListResponse represents the response with books from subjects
type BookListResponse struct {
	Books   []BookWithSubject `json:"books"`
	Total   int               `json:"total"`
	Message string            `json:"message,omitempty"`
}

// BookWithSubject includes subject information with each book
type BookWithSubject struct {
	Book
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	SubjectURL  string    `json:"subject_url"`
}

// PDFDownloadRequest represents a request to download PDFs from selected books
type PDFDownloadRequest struct {
	Books []BookDownloadInfo `json:"books"`
}

// BookDownloadInfo contains the information needed to download a book PDF
type BookDownloadInfo struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	PDFLink string `json:"pdfLink"`
}

// PDFDownloadResponse represents the response from PDF download operation
type PDFDownloadResponse struct {
	JobID     string                   `json:"job_id"`
	Status    string                   `json:"status"`
	Documents []DocumentDownloadResult `json:"documents"`
	Message   string                   `json:"message"`
}

// DocumentDownloadResult represents the result of downloading a single document
type DocumentDownloadResult struct {
	Title      string    `json:"title"`
	DocumentID uuid.UUID `json:"document_id"`
	Status     string    `json:"status"`
	Error      string    `json:"error,omitempty"`
}

package content_discovery

import (
	"github.com/google/uuid"
)

// Book represents a book found from a subject URL
type Book struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	Authors     string `json:"authors,omitempty"`
	Selected    bool   `json:"selected"`
}

// BookListRequest represents a request to list books from subjects
type BookListRequest struct {
	SubjectIDs []uuid.UUID `json:"subject_ids"`
	MaxBooks   int         `json:"max_books,omitempty"` // Default to 10 if not specified
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

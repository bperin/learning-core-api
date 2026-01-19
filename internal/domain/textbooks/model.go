package textbooks

import (
	"time"

	"github.com/google/uuid"
)

// Subject represents a subject from the Open Textbook Library
type Subject struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	SubSubjects []SubSubject `json:"sub_subjects"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// SubSubject represents a sub-subject within a subject
type SubSubject struct {
	ID        uuid.UUID `json:"id"`
	SubjectID uuid.UUID `json:"subject_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScraperResult represents the result of a scraping operation
type ScraperResult struct {
	Subjects  []Subject `json:"subjects"`
	ScrapedAt time.Time `json:"scraped_at"`
}

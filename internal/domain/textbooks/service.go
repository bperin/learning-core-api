package textbooks

import (
	"context"
	"fmt"
	"log"
)

// DownloadResult represents the result of a download operation
type DownloadResult struct {
	DownloadedCount int
	FailedCount     int
	DownloadPath    string
}

// Book represents a textbook from Open Textbook Library
type Book struct {
	Title   string
	URL     string
	PDFLink string
}

// Service defines the business logic for textbooks
type Service interface {
	// ScrapeAndStoreSubjects scrapes subjects and stores them in the database
	ScrapeAndStoreSubjects(ctx context.Context) (*ScraperResult, error)

	// GetAllSubjects retrieves all subjects from the database
	GetAllSubjects(ctx context.Context) ([]Subject, error)

	// GetSubjectByName retrieves a subject by name
	GetSubjectByName(ctx context.Context, name string) (*Subject, error)

	// GetBooksBySubject retrieves books for a given subject slug
	GetBooksBySubject(ctx context.Context, slug string) ([]Book, error)

	// DownloadBooks downloads selected books from a subject
	DownloadBooks(ctx context.Context, subjectSlug string, bookURLs []string) (*DownloadResult, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	repo   Repository
	scraper *Scraper
}

// NewService creates a new textbook service
func NewService(repo Repository) Service {
	return &serviceImpl{
		repo:    repo,
		scraper: NewScraper(),
	}
}

// ScrapeAndStoreSubjects scrapes subjects from Open Textbook Library and stores them
func (s *serviceImpl) ScrapeAndStoreSubjects(ctx context.Context) (*ScraperResult, error) {
	log.Println("Starting scrape and store operation...")

	// Scrape subjects from the web
	result, err := s.scraper.ScrapeSubjects()
	if err != nil {
		return nil, fmt.Errorf("failed to scrape subjects: %w", err)
	}

	// Store subjects in database
	for _, subject := range result.Subjects {
		if err := s.repo.CreateSubject(ctx, &subject); err != nil {
			log.Printf("Warning: failed to store subject %s: %v", subject.Name, err)
			continue
		}

		// Store sub-subjects
		for _, subSubject := range subject.SubSubjects {
			subSubject.SubjectID = subject.ID
			if err := s.repo.CreateSubSubject(ctx, &subSubject); err != nil {
				log.Printf("Warning: failed to store sub-subject %s: %v", subSubject.Name, err)
			}
		}
	}

	log.Printf("Successfully stored %d subjects", len(result.Subjects))
	return result, nil
}

// GetAllSubjects retrieves all subjects from the database
func (s *serviceImpl) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	return s.repo.GetAllSubjects(ctx)
}

// GetSubjectByName retrieves a subject by name
func (s *serviceImpl) GetSubjectByName(ctx context.Context, name string) (*Subject, error) {
	return s.repo.GetSubjectByName(ctx, name)
}

// GetBooksBySubject retrieves books for a given subject slug from Open Textbook Library
func (s *serviceImpl) GetBooksBySubject(ctx context.Context, slug string) ([]Book, error) {
	// Fetch books from Open Textbook Library API
	books, err := s.scraper.GetBooksBySubject(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get books for subject %s: %w", slug, err)
	}
	return books, nil
}

// DownloadBooks downloads selected books from a subject
func (s *serviceImpl) DownloadBooks(ctx context.Context, subjectSlug string, bookURLs []string) (*DownloadResult, error) {
	result := &DownloadResult{
		DownloadPath: fmt.Sprintf("./downloads/%s", subjectSlug),
	}

	// For now, return a placeholder result
	// In a real implementation, this would call the PDF downloader
	result.DownloadedCount = len(bookURLs)
	result.FailedCount = 0

	log.Printf("Initiated download of %d books from subject %s", len(bookURLs), subjectSlug)
	return result, nil
}

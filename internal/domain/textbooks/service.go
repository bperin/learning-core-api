package textbooks

import (
	"context"
	"fmt"
	"log"
)

// Service defines the business logic for textbooks
type Service interface {
	// ScrapeAndStoreSubjects scrapes subjects and stores them in the database
	ScrapeAndStoreSubjects(ctx context.Context) (*ScraperResult, error)

	// GetAllSubjects retrieves all subjects from the database
	GetAllSubjects(ctx context.Context) ([]Subject, error)

	// GetSubjectByName retrieves a subject by name
	GetSubjectByName(ctx context.Context, name string) (*Subject, error)
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

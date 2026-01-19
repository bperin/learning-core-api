package content_discovery

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"

	"learning-core-api/internal/domain/subjects"
)

type Service struct {
	subjectsService subjects.Service
}

func NewService(subjectsService subjects.Service) *Service {
	return &Service{
		subjectsService: subjectsService,
	}
}

// ListBooksFromSubjects fetches books from the provided subject URLs
func (s *Service) ListBooksFromSubjects(ctx context.Context, req BookListRequest) (*BookListResponse, error) {
	if len(req.SubjectIDs) == 0 {
		return &BookListResponse{
			Books:   []BookWithSubject{},
			Total:   0,
			Message: "No subjects provided",
		}, nil
	}

	// Set default max books if not specified
	maxBooks := req.MaxBooks
	if maxBooks <= 0 {
		maxBooks = 10
	}

	// Get all subjects to find the ones requested
	allSubjects, err := s.subjectsService.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}

	// Create a map for quick lookup
	subjectMap := make(map[uuid.UUID]subjects.Subject)
	subSubjectMap := make(map[uuid.UUID]subjects.SubSubject)
	
	for _, subject := range allSubjects {
		subjectMap[subject.ID] = subject
		for _, subSubject := range subject.SubSubjects {
			subSubjectMap[subSubject.ID] = subSubject
		}
	}

	var allBooks []BookWithSubject
	
	// Process each requested subject
	for _, subjectID := range req.SubjectIDs {
		var subjectURL, subjectName string
		
		// Check if it's a main subject or sub-subject
		if subject, exists := subjectMap[subjectID]; exists {
			subjectURL = subject.Url
			subjectName = subject.Name
		} else if subSubject, exists := subSubjectMap[subjectID]; exists {
			subjectURL = subSubject.Url
			subjectName = subSubject.Name
		} else {
			continue // Skip unknown subject IDs
		}

		// Fetch books from this subject
		books, err := s.fetchBooksFromURL(subjectURL)
		if err != nil {
			// Log error but continue with other subjects
			continue
		}

		// Add subject info to each book
		for _, book := range books {
			if len(allBooks) >= maxBooks {
				break
			}
			
			allBooks = append(allBooks, BookWithSubject{
				Book:        book,
				SubjectID:   subjectID,
				SubjectName: subjectName,
				SubjectURL:  subjectURL,
			})
		}

		if len(allBooks) >= maxBooks {
			break
		}
	}

	return &BookListResponse{
		Books: allBooks,
		Total: len(allBooks),
		Message: fmt.Sprintf("Found %d books from %d subjects", len(allBooks), len(req.SubjectIDs)),
	}, nil
}

// fetchBooksFromURL scrapes books from a subject URL (similar to pdf-downloader logic)
func (s *Service) fetchBooksFromURL(subjectURL string) ([]Book, error) {
	resp, err := http.Get(subjectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch subject page: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var books []Book

	// Look for book links (similar to pdf-downloader pattern)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Look for textbook links (adapt this pattern based on the actual site structure)
		if strings.Contains(href, "/textbooks/") || strings.Contains(href, "/book/") {
			title := strings.TrimSpace(s.Text())
			if title == "" {
				// Try to get title from parent or child elements
				title = strings.TrimSpace(s.Parent().Text())
			}
			
			if title != "" && len(title) > 3 {
				// Clean up the title
				title = cleanTitle(title)
				
				// Make sure we have a full URL
				fullURL := href
				if strings.HasPrefix(href, "/") {
					fullURL = "https://open.umn.edu" + href
				}

				books = append(books, Book{
					Title:    title,
					URL:      fullURL,
					Selected: false,
				})
			}
		}
	})

	// Remove duplicates
	books = removeDuplicateBooks(books)

	// Limit to reasonable number per subject
	if len(books) > 20 {
		books = books[:20]
	}

	return books, nil
}

// cleanTitle removes extra whitespace and common prefixes/suffixes
func cleanTitle(title string) string {
	// Remove extra whitespace
	title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)
	
	// Remove common patterns
	title = strings.TrimPrefix(title, "View ")
	title = strings.TrimSuffix(title, " - Open Textbook Library")
	
	return title
}

// removeDuplicateBooks removes duplicate books based on title
func removeDuplicateBooks(books []Book) []Book {
	seen := make(map[string]bool)
	var result []Book
	
	for _, book := range books {
		if !seen[book.Title] {
			seen[book.Title] = true
			result = append(result, book)
		}
	}
	
	return result
}

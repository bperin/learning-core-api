package content_discovery

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/subjects"
)

// Atom feed structures for parsing Open Textbook Library
type AtomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	ID    string   `xml:"id"`
	Title string   `xml:"title"`
	Link  AtomLink `xml:"link"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

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
	for _, subjectIDStr := range req.SubjectIDs {
		// Parse string UUID to uuid.UUID
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			continue // Skip invalid UUIDs
		}

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

// fetchBooksFromURL fetches books from a subject URL using the Atom XML feed
// and extracts PDF links for each book
func (s *Service) fetchBooksFromURL(subjectURL string) ([]Book, error) {
	// Create request with Accept header for Atom feed
	req, err := http.NewRequest("GET", subjectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/atom+xml")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch subject page: status %d", resp.StatusCode)
	}

	// Parse the Atom feed
	var feed AtomFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to parse Atom feed: %w", err)
	}

	var books []Book
	for _, entry := range feed.Entries {
		if entry.Link.Href != "" && entry.Title != "" {
			// Fetch PDF link from the book's detail page
			pdfLink, err := s.getPDFLinkFromBookPage(entry.Link.Href)
			if err != nil || pdfLink == "" {
				// Skip books without PDF links
				continue
			}

			book := Book{
				Title:    strings.TrimSpace(entry.Title),
				URL:      entry.Link.Href,
				PDFLink:  pdfLink,
				Selected: false,
			}

			books = append(books, book)
		}
	}

	return books, nil
}

// getPDFLinkFromBookPage extracts the PDF download link from a book detail page
// Returns the actual PDF download URL by following the format page
func (s *Service) getPDFLinkFromBookPage(bookURL string) (string, error) {
	resp, err := http.Get(bookURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch book page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	htmlContent := string(body)

	// Look for the PDF format link pattern: /opentextbooks/formats/XXX
	formatsIdx := strings.Index(htmlContent, "/opentextbooks/formats/")
	if formatsIdx != -1 {
		urlStart := formatsIdx
		urlEnd := strings.IndexAny(htmlContent[urlStart:], "\")")
		if urlEnd != -1 {
			formatPageURL := "https://open.umn.edu" + htmlContent[urlStart:urlStart+urlEnd]
			// Now fetch the format page to get the actual PDF download link
			return s.getDownloadLinkFromFormatPage(formatPageURL)
		}
	}

	return "", nil
}

// getDownloadLinkFromFormatPage fetches the format page and extracts the actual PDF download link
func (s *Service) getDownloadLinkFromFormatPage(formatPageURL string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // Allow redirects
		},
	}

	resp, err := client.Get(formatPageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch format page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read format page: %w", err)
	}

	htmlContent := string(body)
	finalURL := resp.Request.URL.String()

	// Strategy 1: Look for direct .pdf links
	if link := s.findPDFLink(htmlContent, ".pdf", finalURL); link != "" {
		return link, nil
	}

	// Strategy 2: Look for download buttons
	if link := s.findPDFLink(htmlContent, "download", finalURL); link != "" {
		return link, nil
	}

	// Strategy 3: Look for "Digital PDF" text
	if link := s.findPDFLink(htmlContent, "digital pdf", finalURL); link != "" {
		return link, nil
	}

	// Strategy 4: Look for type=pdf parameter
	if link := s.findPDFLink(htmlContent, "type=pdf", finalURL); link != "" {
		return link, nil
	}

	// Strategy 5: Search all hrefs for pdf/download keywords
	return s.findAnyPDFHref(htmlContent, finalURL), nil
}

// findPDFLink searches for a pattern and extracts the href before it
func (s *Service) findPDFLink(htmlContent, pattern, baseURL string) string {
	idx := strings.Index(strings.ToLower(htmlContent), strings.ToLower(pattern))
	if idx == -1 {
		return ""
	}

	searchStart := idx
	if searchStart > 200 {
		searchStart = idx - 200
	} else {
		searchStart = 0
	}

	searchText := htmlContent[searchStart:idx]
	hrefIdx := strings.LastIndex(searchText, "href=\"")
	if hrefIdx == -1 {
		return ""
	}

	hrefStart := searchStart + hrefIdx + 6
	hrefEnd := strings.Index(htmlContent[hrefStart:], "\"")
	if hrefEnd == -1 {
		return ""
	}

	link := htmlContent[hrefStart : hrefStart+hrefEnd]
	return s.makeAbsoluteURL(link, baseURL)
}

// findAnyPDFHref searches all hrefs for pdf or download keywords
func (s *Service) findAnyPDFHref(htmlContent, baseURL string) string {
	lowerContent := strings.ToLower(htmlContent)
	idx := strings.Index(lowerContent, "href=\"")

	for idx != -1 {
		hrefStart := idx + 6
		hrefEnd := strings.Index(htmlContent[hrefStart:], "\"")
		if hrefEnd == -1 {
			break
		}

		link := htmlContent[hrefStart : hrefStart+hrefEnd]
		lowerLink := strings.ToLower(link)
		if strings.Contains(lowerLink, "pdf") || strings.Contains(lowerLink, "download") {
			return s.makeAbsoluteURL(link, baseURL)
		}

		// Find next href
		nextIdx := strings.Index(lowerContent[hrefStart+hrefEnd:], "href=\"")
		if nextIdx == -1 {
			break
		}
		idx = hrefStart + hrefEnd + nextIdx
	}

	return ""
}

// makeAbsoluteURL converts a relative URL to absolute
func (s *Service) makeAbsoluteURL(link, baseURL string) string {
	if strings.HasPrefix(link, "http") {
		return link
	}
	if strings.HasPrefix(link, "/") {
		parts := strings.Split(baseURL, "/")
		if len(parts) >= 3 {
			domain := strings.Join(parts[:3], "/")
			return domain + link
		}
	}
	return baseURL + link
}

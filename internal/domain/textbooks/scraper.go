package textbooks

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
)

const baseURL = "https://open.umn.edu/opentextbooks/subjects"

// Scraper handles scraping of Open Textbook Library subjects
type Scraper struct {
	client *http.Client
}

// NewScraper creates a new textbook scraper
func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// ScrapeSubjects scrapes all subjects and their sub-subjects from the Open Textbook Library
func (s *Scraper) ScrapeSubjects() (*ScraperResult, error) {
	log.Println("Starting to scrape Open Textbook Library subjects...")

	subjects, err := s.scrapeSubjectsPage()
	if err != nil {
		return nil, fmt.Errorf("failed to scrape subjects: %w", err)
	}

	log.Printf("Successfully scraped %d subjects", len(subjects))

	return &ScraperResult{
		Subjects:  subjects,
		ScrapedAt: time.Now().UTC(),
	}, nil
}

func (s *Scraper) scrapeSubjectsPage() ([]Subject, error) {
	resp, err := s.client.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var subjects []Subject
	seen := make(map[string]bool)

	// Find all subject links
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Filter for subject URLs
		if !strings.Contains(href, "/opentextbooks/subjects/") {
			return
		}

		// Skip the main subjects page link
		if href == "/opentextbooks/subjects" || href == "/opentextbooks/subjects/" {
			return
		}

		// Skip if we've already seen this URL
		if seen[href] {
			return
		}

		text := strings.TrimSpace(sel.Text())
		if text == "" {
			return
		}

		seen[href] = true

		subject := Subject{
			ID:          uuid.New(),
			Name:        text,
			URL:         href,
			SubSubjects: []SubSubject{},
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		// Scrape sub-subjects for this subject
		subSubjects, err := s.scrapeSubSubjectsPage(href)
		if err != nil {
			log.Printf("Warning: failed to scrape sub-subjects for %s: %v", text, err)
		} else {
			subject.SubSubjects = subSubjects
		}

		subjects = append(subjects, subject)
	})

	return subjects, nil
}

func (s *Scraper) scrapeSubSubjectsPage(subjectURL string) ([]SubSubject, error) {
	// Ensure the URL is absolute
	fullURL := subjectURL
	if !strings.HasPrefix(fullURL, "http") {
		fullURL = "https://open.umn.edu" + subjectURL
	}

	resp, err := s.client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var subSubjects []SubSubject
	seen := make(map[string]bool)

	// Look for sub-subject links within the subject page
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Filter for sub-subject URLs
		if !strings.Contains(href, "/opentextbooks/subjects/") {
			return
		}

		// Skip if it's the same as the parent subject
		if href == subjectURL {
			return
		}

		// Skip if we've already seen this URL
		if seen[href] {
			return
		}

		text := strings.TrimSpace(sel.Text())
		if text == "" || text == "Subjects" {
			return
		}

		seen[href] = true
		subSubjects = append(subSubjects, SubSubject{
			ID:        uuid.New(),
			Name:      text,
			URL:       href,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	})

	return subSubjects, nil
}

// GetBooksBySubject fetches books for a given subject slug from the Atom feed
func (s *Scraper) GetBooksBySubject(slug string) ([]Book, error) {
	var books []Book

	// Fetch the Atom feed for the subject
	subjectURL := fmt.Sprintf("%s/%s", baseURL, slug)
	req, err := http.NewRequest("GET", subjectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/atom+xml")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subject: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML to extract book links (fallback if Atom doesn't work)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Find all book links
	doc.Find("a[href*='/opentextbooks/textbooks/']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		title := strings.TrimSpace(s.Text())
		if title == "" {
			return
		}

		// Make URL absolute if needed
		if !strings.HasPrefix(href, "http") {
			href = "https://open.umn.edu" + href
		}

		books = append(books, Book{
			Title: title,
			URL:   href,
		})
	})

	return books, nil
}

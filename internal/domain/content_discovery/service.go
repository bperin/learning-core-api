package content_discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/domain/generation"
	"learning-core-api/internal/domain/subjects"
	"learning-core-api/internal/gcp"
	"learning-core-api/internal/infra/progress"

	"github.com/google/uuid"
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
	subjectsService   subjects.Service
	documentsService  documents.Service
	gcsService        *gcp.GCSService
	fileService       *gcp.FileService
	generationService *generation.Service
}

func NewService(subjectsService subjects.Service, documentsService documents.Service, gcsService *gcp.GCSService, fileService *gcp.FileService, generationService *generation.Service) *Service {
	if subjectsService == nil {
		panic("subjectsService is required")
	}
	if documentsService == nil {
		panic("documentsService is required")
	}
	if gcsService == nil {
		panic("gcsService is required")
	}
	if fileService == nil {
		panic("fileService is required")
	}
	return &Service{
		subjectsService:   subjectsService,
		documentsService:  documentsService,
		gcsService:        gcsService,
		fileService:       fileService,
		generationService: generationService,
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
		Books:   allBooks,
		Total:   len(allBooks),
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

// DownloadPDFs downloads PDFs from the provided book links and processes them
func (s *Service) DownloadPDFs(ctx context.Context, req PDFDownloadRequest, userIDStr string) (*PDFDownloadResponse, error) {
	if len(req.Books) == 0 {
		log.Printf("[PDF_DOWNLOAD] No books provided in request")
		return &PDFDownloadResponse{
			Status:  "error",
			Message: "No books provided",
		}, nil
	}

	jobID := uuid.New().String()
	log.Printf("[PDF_DOWNLOAD] Starting job %s with %d books", jobID, len(req.Books))

	// Initialize progress tracking
	tracker := progress.GetTracker()
	tracker.StartJob(jobID, "pdf_download")
	tracker.UpdateProgress(jobID, "started", fmt.Sprintf("Starting PDF download for %d books", len(req.Books)), 0, nil)

	// Parse user ID from string
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[PDF_DOWNLOAD] Failed to parse user ID: %v", err)
		tracker.FailJob(jobID, fmt.Sprintf("Invalid user ID: %v", err))
		return &PDFDownloadResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid user ID: %v", err),
		}, err
	}

	results := make([]DocumentDownloadResult, len(req.Books))
	var wg sync.WaitGroup
	var mu sync.Mutex
	completedCount := 0

	// Process books in parallel with goroutines
	for i, book := range req.Books {
		wg.Add(1)
		go func(index int, b BookDownloadInfo) {
			defer wg.Done()

			log.Printf("[PDF_DOWNLOAD] [Job:%s] [Book:%d] Starting download: %s", jobID, index+1, b.Title)

			result := DocumentDownloadResult{
				Title:  b.Title,
				Status: "processing",
			}

			documentID, err := s.downloadAndProcessPDF(ctx, b, userID)
			if err != nil {
				log.Printf("[PDF_DOWNLOAD] [Job:%s] [Book:%d] FAILED: %s - Error: %v", jobID, index+1, b.Title, err)
				result.Status = "failed"
				result.Error = err.Error()
				tracker.UpdateProgressWithError(jobID, "processing", fmt.Sprintf("Failed to download %s", b.Title), err.Error(), 0)
			} else {
				log.Printf("[PDF_DOWNLOAD] [Job:%s] [Book:%d] SUCCESS: %s - DocumentID: %s", jobID, index+1, b.Title, documentID)
				result.DocumentID = documentID
				result.Status = "success"

				mu.Lock()
				completedCount++
				progress := (completedCount * 100) / len(req.Books)
				mu.Unlock()

				tracker.UpdateProgress(jobID, "processing", fmt.Sprintf("Downloaded %d/%d books", completedCount, len(req.Books)), progress, nil)
			}

			results[index] = result
		}(i, book)
	}

	// Wait for all downloads to complete
	log.Printf("[PDF_DOWNLOAD] [Job:%s] Waiting for all %d downloads to complete...", jobID, len(req.Books))
	wg.Wait()

	// Count successful downloads
	successCount := 0
	failedCount := 0
	for _, result := range results {
		if result.Status == "completed" {
			successCount++
		} else {
			failedCount++
		}
	}

	log.Printf("[PDF_DOWNLOAD] [Job:%s] COMPLETED - Success: %d, Failed: %d", jobID, successCount, failedCount)

	return &PDFDownloadResponse{
		JobID:     jobID,
		Status:    "completed",
		Documents: results,
		Message:   fmt.Sprintf("Processed %d books, %d successful downloads", len(req.Books), successCount),
	}, nil
}

// downloadAndProcessPDF downloads a single PDF and processes it through the pipeline
func (s *Service) downloadAndProcessPDF(ctx context.Context, book BookDownloadInfo, userID uuid.UUID) (uuid.UUID, error) {
	log.Printf("[PDF_PROCESS] Starting processing for: %s", book.Title)
	log.Printf("[PDF_PROCESS] PDF URL: %s", book.PDFLink)

	// Validate PDF URL before attempting download
	if book.PDFLink == "" {
		return uuid.Nil, fmt.Errorf("PDF link is empty")
	}

	if !strings.HasPrefix(book.PDFLink, "http://") && !strings.HasPrefix(book.PDFLink, "https://") {
		return uuid.Nil, fmt.Errorf("invalid PDF URL scheme: %s", book.PDFLink)
	}

	// Step 1: Download PDF from URL
	log.Printf("[PDF_PROCESS] Step 1: Downloading PDF from URL...")
	resp, err := http.Get(book.PDFLink)
	if err != nil {
		log.Printf("[PDF_PROCESS] Step 1 FAILED: HTTP request failed - %v", err)
		return uuid.Nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[PDF_PROCESS] Step 1 FAILED: HTTP status %d", resp.StatusCode)
		return uuid.Nil, fmt.Errorf("failed to download PDF: status %d", resp.StatusCode)
	}

	// Buffer the response body so we can use it multiple times
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[PDF_PROCESS] Step 1 FAILED: Failed to read response body - %v", err)
		return uuid.Nil, fmt.Errorf("failed to read PDF data: %w", err)
	}

	if len(pdfData) == 0 {
		log.Printf("[PDF_PROCESS] Step 1 FAILED: PDF data is empty")
		return uuid.Nil, fmt.Errorf("PDF data is empty")
	}

	log.Printf("[PDF_PROCESS] Step 1 SUCCESS: PDF downloaded, Size: %d bytes", len(pdfData))

	// Step 2: Create document record in database
	log.Printf("[PDF_PROCESS] Step 2: Creating document record in database...")
	filename := fmt.Sprintf("%s.pdf", sanitizeFilename(book.Title))
	title := book.Title
	mimeType := "application/pdf"

	// Create document using documents service
	createReq := documents.CreateDocumentRequest{
		Filename:  filename,
		Title:     &title,
		MimeType:  &mimeType,
		RagStatus: documents.RagStatusPending,
		UserID:    userID,
	}

	doc, err := s.documentsService.CreateDocument(ctx, createReq)
	if err != nil {
		log.Printf("[PDF_PROCESS] Step 2 FAILED: Database insert failed - %v", err)
		return uuid.Nil, fmt.Errorf("failed to create document record: %w", err)
	}

	log.Printf("[PDF_PROCESS] Step 2 SUCCESS: Document created with ID: %s", doc.ID)

	// Step 3: Upload PDF to GCS first (following test pattern)
	objectName := fmt.Sprintf("documents/%s/%s", doc.ID, filename)
	log.Printf("[PDF_PROCESS] Step 3: Uploading to GCS with object name: %s", objectName)

	// Upload to GCS using GCSService with buffered data
	_, err = s.gcsService.UploadFile(ctx, objectName, "application/pdf", bytes.NewReader(pdfData))
	if err != nil {
		log.Printf("[PDF_PROCESS] Step 3 FAILED: GCS upload failed - %v", err)
		return uuid.Nil, fmt.Errorf("failed to upload to GCS: %w", err)
	}

	log.Printf("[PDF_PROCESS] Step 3 SUCCESS: File uploaded to GCS")

	// Step 4: Upload to File Search Store using the GCS object name (like line 82 in test)
	log.Printf("[PDF_PROCESS] Step 4: Uploading to File Search Store...")
	result, err := s.fileService.UploadToFileSearchStore(ctx, objectName, book.Title, "application/pdf")
	if err != nil {
		log.Printf("[PDF_PROCESS] Step 4 FAILED: File Search Store upload failed - %v", err)
		// Cleanup: Delete from GCS if File Search Store upload fails
		log.Printf("[PDF_PROCESS] Cleaning up: Deleting %s from GCS...", objectName)
		if cleanupErr := s.gcsService.DeleteObject(ctx, objectName); cleanupErr != nil {
			log.Printf("[PDF_PROCESS] Cleanup FAILED: Failed to delete %s from GCS: %v", objectName, cleanupErr)
		}
		return uuid.Nil, fmt.Errorf("failed to upload to file search store: %w", err)
	}

	log.Printf("[PDF_PROCESS] Step 4 SUCCESS: File uploaded to File Search Store")
	if result != nil && result.Operation != nil {
		log.Printf("[PDF_PROCESS] File Search Store operation started: %s", result.Operation.Name)
	}

	// Step 5: Trigger classification generation asynchronously (if generation service is available)
	if s.generationService != nil {
		go func() {
			log.Printf("[PDF_CLASSIFICATION] Starting classification generation for document: %s", doc.ID)
			err := s.generateClassificationArtifacts(context.Background(), doc.ID, book.Title)
			if err != nil {
				log.Printf("[PDF_CLASSIFICATION] FAILED for document %s: %v", doc.ID, err)
			} else {
				log.Printf("[PDF_CLASSIFICATION] SUCCESS for document %s", doc.ID)
			}
		}()
	} else {
		log.Printf("[PDF_CLASSIFICATION] Skipping classification generation - generation service not available")
	}

	log.Printf("[PDF_PROCESS] COMPLETED: All steps successful for %s (DocumentID: %s)", book.Title, doc.ID)
	return doc.ID, nil
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(title string) string {
	// Replace invalid filename characters
	filename := strings.ReplaceAll(title, "/", "-")
	filename = strings.ReplaceAll(filename, "\\", "-")
	filename = strings.ReplaceAll(filename, ":", "-")
	filename = strings.ReplaceAll(filename, "*", "-")
	filename = strings.ReplaceAll(filename, "?", "-")
	filename = strings.ReplaceAll(filename, "\"", "-")
	filename = strings.ReplaceAll(filename, "<", "-")
	filename = strings.ReplaceAll(filename, ">", "-")
	filename = strings.ReplaceAll(filename, "|", "-")

	// Limit length
	if len(filename) > 100 {
		filename = filename[:100]
	}

	return filename
}

// generateClassificationArtifacts triggers classification generation for a document using the file search store
func (s *Service) generateClassificationArtifacts(ctx context.Context, documentID uuid.UUID, title string) error {
	log.Printf("[PDF_CLASSIFICATION] Generating classification artifacts for document: %s (%s)", documentID, title)

	// Create file search tool config with the document's file search store reference
	fileSearchConfig := map[string]interface{}{
		"store_names": []string{s.fileService.StoreName()},
	}

	fileSearchConfigJSON, err := json.Marshal(fileSearchConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal file search config: %w", err)
	}

	// Create generation request for classification
	generateReq := generation.GenerateRequest{
		UserID: uuid.New(), // TODO: Get actual user ID from context
		Target: generation.Target{
			DocumentID: &documentID,
		},
		Instructions: generation.Instructions{
			GenerationType: "classification", // This should reference a classification prompt template in DB
			PromptVersion:  0,                // Use latest version
			Variables: map[string]interface{}{
				"document_title": title,
				"document_id":    documentID.String(),
			},
		},
		Output: generation.OutputConfig{
			GenerationType: "classification", // This should reference a classification schema template in DB
			SchemaVersion:  0,                // Use latest version
			Format:         "json",
		},
		Tools: []generation.ToolConfig{
			{
				Type:   "file_search",
				Config: fileSearchConfigJSON,
			},
		},
		ModelConfigID: uuid.Nil, // Use active model config
	}

	log.Printf("[PDF_CLASSIFICATION] Calling generation service for document: %s", documentID)
	response, err := s.generationService.Generate(ctx, generateReq)
	if err != nil {
		return fmt.Errorf("failed to generate classification: %w", err)
	}

	log.Printf("[PDF_CLASSIFICATION] Classification generated successfully for document: %s, ArtifactID: %s", documentID, response.ArtifactID)

	// The generation service automatically saves the artifact to the database
	// The artifact contains the classification results and can be retrieved later for question generation

	return nil
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

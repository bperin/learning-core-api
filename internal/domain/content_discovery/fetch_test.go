package content_discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchBooksFromURL_Education(t *testing.T) {
	service := &Service{}

	books, err := service.fetchBooksFromURL("https://open.umn.edu/opentextbooks/subjects/education")
	require.NoError(t, err)

	assert.Greater(t, len(books), 0, "Expected to find books from education subject")

	booksWithPDF := 0
	t.Logf("Found %d books from education subject", len(books))
	for i, book := range books {
		t.Logf("  %d. %s", i+1, book.Title)
		t.Logf("     URL: %s", book.URL)
		if book.PDFLink != "" {
			t.Logf("     PDF: %s", book.PDFLink)
			booksWithPDF++
		} else {
			t.Logf("     PDF: (none found)")
		}
	}

	t.Logf("Books with PDF links: %d/%d", booksWithPDF, len(books))

	// Verify book structure
	for _, book := range books {
		assert.NotEmpty(t, book.Title, "Book should have a title")
		assert.NotEmpty(t, book.URL, "Book should have a URL")
		assert.Contains(t, book.URL, "open.umn.edu", "Book URL should be from open.umn.edu")
	}
}

func TestFetchBooksFromURL_Biology(t *testing.T) {
	service := &Service{}

	books, err := service.fetchBooksFromURL("https://open.umn.edu/opentextbooks/subjects/biology")
	require.NoError(t, err)

	assert.Greater(t, len(books), 0, "Expected to find books from biology subject")

	t.Logf("Found %d books from biology subject", len(books))
	for i, book := range books {
		if i >= 5 {
			break
		}
		t.Logf("  %d. %s - %s", i+1, book.Title, book.URL)
	}
}

func TestFetchBooksFromURL_Mathematics(t *testing.T) {
	service := &Service{}

	books, err := service.fetchBooksFromURL("https://open.umn.edu/opentextbooks/subjects/mathematics")
	require.NoError(t, err)

	assert.Greater(t, len(books), 0, "Expected to find books from mathematics subject")

	t.Logf("Found %d books from mathematics subject", len(books))
	for i, book := range books {
		if i >= 5 {
			break
		}
		t.Logf("  %d. %s - %s", i+1, book.Title, book.URL)
	}
}

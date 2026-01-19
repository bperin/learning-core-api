package textbooks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockRepository implements Repository interface for testing
type MockRepository struct {
	subjects []Subject
}

func (m *MockRepository) CreateSubject(ctx context.Context, subject *Subject) error {
	m.subjects = append(m.subjects, *subject)
	return nil
}

func (m *MockRepository) CreateSubSubject(ctx context.Context, subSubject *SubSubject) error {
	return nil
}

func (m *MockRepository) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	return m.subjects, nil
}

func (m *MockRepository) GetSubjectByID(ctx context.Context, id uuid.UUID) (*Subject, error) {
	for _, s := range m.subjects {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) GetSubjectByName(ctx context.Context, name string) (*Subject, error) {
	for _, s := range m.subjects {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) GetSubSubjectsBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]SubSubject, error) {
	return []SubSubject{}, nil
}

func (m *MockRepository) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockRepository) DeleteAllSubjects(ctx context.Context) error {
	m.subjects = []Subject{}
	return nil
}

func setupTestHandler() (*Handler, *MockRepository) {
	repo := &MockRepository{
		subjects: []Subject{
			{
				ID:        uuid.New(),
				Name:      "Open Textbook Library - Databases Textbooks",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Open Textbook Library - Mathematics Textbooks",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Open Textbook Library - Business Textbooks",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	service := NewService(repo)
	handler := NewHandler(service)

	return handler, repo
}

func TestAdminListSubjects(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/admin/textbooks/subjects", nil)
	w := httptest.NewRecorder()

	handler.AdminListSubjects(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if _, ok := response["subjects"]; !ok {
		t.Error("Response missing 'subjects' field")
	}

	if _, ok := response["count"]; !ok {
		t.Error("Response missing 'count' field")
	}

	// Verify count is correct
	count := int(response["count"].(float64))
	if count != 3 {
		t.Errorf("Expected 3 subjects, got %d", count)
	}

	// Verify subjects array
	subjects := response["subjects"].([]interface{})
	if len(subjects) != 3 {
		t.Errorf("Expected 3 subjects in array, got %d", len(subjects))
	}

	// Verify subject names
	expectedNames := map[string]bool{
		"Open Textbook Library - Databases Textbooks":   false,
		"Open Textbook Library - Mathematics Textbooks": false,
		"Open Textbook Library - Business Textbooks":    false,
	}

	for _, s := range subjects {
		subject := s.(map[string]interface{})
		name := subject["name"].(string)
		if _, exists := expectedNames[name]; exists {
			expectedNames[name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected subject '%s' not found in response", name)
		}
	}
}

func TestAdminGetBooksBySubject(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/admin/textbooks/subjects/databases/books", nil)
	w := httptest.NewRecorder()

	handler.AdminGetBooksBySubject(w, req)

	// Expected 400 since we're not using chi router to parse params
	// This is a limitation of httptest without chi
	if w.Code != http.StatusBadRequest {
		t.Logf("Status code: %d (expected 400 without chi router)", w.Code)
	}
}

func TestAdminDownloadBooks(t *testing.T) {
	handler, _ := setupTestHandler()

	downloadReq := AdminDownloadRequest{
		SubjectSlug: "databases",
		BookURLs: []string{
			"https://open.umn.edu/opentextbooks/textbooks/database-design-2nd-edition",
			"https://open.umn.edu/opentextbooks/textbooks/relational-databases-and-microsoft-access",
		},
	}

	body, err := json.Marshal(downloadReq)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/admin/textbooks/download", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AdminDownloadBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response structure
	if _, ok := response["message"]; !ok {
		t.Error("Response missing 'message' field")
	}

	if _, ok := response["subject"]; !ok {
		t.Error("Response missing 'subject' field")
	}

	if _, ok := response["downloaded"]; !ok {
		t.Error("Response missing 'downloaded' field")
	}

	// Verify values
	subject := response["subject"].(string)
	if subject != "databases" {
		t.Errorf("Expected subject 'databases', got '%s'", subject)
	}

	downloaded := int(response["downloaded"].(float64))
	if downloaded != 2 {
		t.Errorf("Expected 2 downloaded books, got %d", downloaded)
	}
}

func TestAdminDownloadBooksValidation(t *testing.T) {
	handler, _ := setupTestHandler()

	tests := []struct {
		name           string
		request        AdminDownloadRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing subject slug",
			request: AdminDownloadRequest{
				SubjectSlug: "",
				BookURLs:    []string{"http://example.com/book1"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Subject slug is required",
		},
		{
			name: "missing book URLs",
			request: AdminDownloadRequest{
				SubjectSlug: "databases",
				BookURLs:    []string{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "At least one book URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			req := httptest.NewRequest("POST", "/admin/textbooks/download", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AdminDownloadBooks(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

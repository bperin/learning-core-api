package content_discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/subjects"
)

// MockSubjectsService is a mock implementation of subjects.Service
type MockSubjectsService struct {
	mock.Mock
}

func (m *MockSubjectsService) ListAll(ctx context.Context) ([]subjects.Subject, error) {
	args := m.Called(ctx)
	return args.Get(0).([]subjects.Subject), args.Error(1)
}

func setupContentDiscoveryHandler(t *testing.T) (*Handler, *MockSubjectsService) {
	t.Helper()

	mockSubjectsService := &MockSubjectsService{}
	service := NewService(mockSubjectsService)
	handler := NewHandler(service)

	return handler, mockSubjectsService
}

func TestHandler_ListBooks_Success(t *testing.T) {
	handler, mockSubjectsService := setupContentDiscoveryHandler(t)

	// Setup mock data
	subjectID := uuid.New()
	mockSubjects := []subjects.Subject{
		{
			ID:   subjectID,
			Name: "Open Textbook Library - Education Textbooks",
			Url:  "https://open.umn.edu/opentextbooks/subjects/education",
		},
	}

	mockSubjectsService.On("ListAll", mock.Anything).Return(mockSubjects, nil)

	// Create request
	reqBody := BookListRequest{
		SubjectIDs: []uuid.UUID{subjectID},
		MaxBooks:   5,
	}
	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response BookListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.GreaterOrEqual(t, response.Total, 0)
	assert.NotEmpty(t, response.Message)
	assert.IsType(t, []BookWithSubject{}, response.Books)

	mockSubjectsService.AssertExpectations(t)
}

func TestHandler_ListBooks_EmptySubjects(t *testing.T) {
	handler, mockSubjectsService := setupContentDiscoveryHandler(t)

	mockSubjectsService.On("ListAll", mock.Anything).Return([]subjects.Subject{}, nil)

	// Create request with empty subject IDs
	reqBody := BookListRequest{
		SubjectIDs: []uuid.UUID{},
		MaxBooks:   5,
	}
	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response BookListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, response.Total)
	assert.Equal(t, "No subjects provided", response.Message)
	assert.Empty(t, response.Books)

	mockSubjectsService.AssertExpectations(t)
}

func TestHandler_ListBooks_InvalidJSON(t *testing.T) {
	handler, _ := setupContentDiscoveryHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_ListBooks_SubjectsServiceError(t *testing.T) {
	handler, mockSubjectsService := setupContentDiscoveryHandler(t)

	subjectID := uuid.New()
	mockSubjectsService.On("ListAll", mock.Anything).Return([]subjects.Subject{}, assert.AnError)

	reqBody := BookListRequest{
		SubjectIDs: []uuid.UUID{subjectID},
		MaxBooks:   5,
	}
	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockSubjectsService.AssertExpectations(t)
}

func TestHandler_ListBooks_DefaultMaxBooks(t *testing.T) {
	handler, mockSubjectsService := setupContentDiscoveryHandler(t)

	subjectID := uuid.New()
	mockSubjects := []subjects.Subject{
		{
			ID:   subjectID,
			Name: "Test Subject",
			Url:  "https://example.com/test",
		},
	}

	mockSubjectsService.On("ListAll", mock.Anything).Return(mockSubjects, nil)

	// Request without MaxBooks specified (should default to 10)
	reqBody := BookListRequest{
		SubjectIDs: []uuid.UUID{subjectID},
	}
	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response BookListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should handle the request successfully even without MaxBooks
	assert.GreaterOrEqual(t, response.Total, 0)

	mockSubjectsService.AssertExpectations(t)
}

func TestHandler_ListBooks_WithSubSubjects(t *testing.T) {
	handler, mockSubjectsService := setupContentDiscoveryHandler(t)

	subjectID := uuid.New()
	subSubjectID := uuid.New()

	mockSubjects := []subjects.Subject{
		{
			ID:   subjectID,
			Name: "Main Subject",
			Url:  "https://example.com/main",
			SubSubjects: []subjects.SubSubject{
				{
					ID:        subSubjectID,
					SubjectID: subjectID,
					Name:      "Sub Subject",
					Url:       "https://example.com/sub",
				},
			},
		},
	}

	mockSubjectsService.On("ListAll", mock.Anything).Return(mockSubjects, nil)

	// Request books from sub-subject
	reqBody := BookListRequest{
		SubjectIDs: []uuid.UUID{subSubjectID},
		MaxBooks:   5,
	}
	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/content-discovery/books", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response BookListResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Should handle sub-subjects correctly
	assert.GreaterOrEqual(t, response.Total, 0)

	mockSubjectsService.AssertExpectations(t)
}

package subjects

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
)

func setupSubjectsHandler(t *testing.T) (*Handler, func()) {
	t.Helper()

	tx, cleanup := testutil.NewTestTx(t)
	queries := store.New(tx)

	repo := NewRepository(queries)
	service := NewService(repo)
	handler := NewHandler(service)

	return handler, cleanup
}

func TestHandler_ListAll(t *testing.T) {
	handler, cleanup := setupSubjectsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/subjects", nil)
	w := httptest.NewRecorder()

	handler.ListAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []Subject
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Should have seeded subjects
	assert.GreaterOrEqual(t, len(result), 1)
}

func TestHandler_ListForSelection(t *testing.T) {
	handler, cleanup := setupSubjectsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/subjects/for-selection", nil)
	w := httptest.NewRecorder()

	handler.ListForSelection(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []SubjectForSelection
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Should have seeded subjects
	assert.GreaterOrEqual(t, len(result), 1)

	// Verify the structure includes required fields
	if len(result) > 0 {
		subject := result[0]
		assert.NotEmpty(t, subject.ID)
		assert.NotEmpty(t, subject.DisplayName)
		assert.NotEmpty(t, subject.FullName)
		assert.NotEmpty(t, subject.URL)

		// Display name should be different from full name for seeded data
		// (assuming seeded data follows "Open Textbook Library - Subject Textbooks" pattern)
		if subject.FullName != subject.DisplayName {
			assert.NotContains(t, subject.DisplayName, "Open Textbook Library")
			assert.NotContains(t, subject.DisplayName, " Textbooks")
		}
	}
}

func TestExtractDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Standard textbook format",
			fullName: "Open Textbook Library - Education Textbooks",
			expected: "Education",
		},
		{
			name:     "Another subject format",
			fullName: "Open Textbook Library - Mathematics Textbooks",
			expected: "Mathematics",
		},
		{
			name:     "Without textbooks suffix",
			fullName: "Open Textbook Library - Science",
			expected: "Science",
		},
		{
			name:     "No dash pattern",
			fullName: "Computer Science",
			expected: "Computer Science",
		},
		{
			name:     "Empty string",
			fullName: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDisplayName(tt.fullName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandler_ListForSelection_EmptyDatabase(t *testing.T) {
	handler, cleanup := setupSubjectsHandler(t)
	defer cleanup()

	// Truncate subjects table to test empty case
	tx, txCleanup := testutil.NewTestTx(t)
	defer txCleanup()

	queries := store.New(tx)
	// Clear any seeded data for this test
	_, err := tx.Exec("TRUNCATE TABLE sub_subjects, subjects RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	repo := NewRepository(queries)
	service := NewService(repo)
	emptyHandler := NewHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/subjects/for-selection", nil)
	w := httptest.NewRecorder()

	emptyHandler.ListForSelection(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []SubjectForSelection
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, 0, len(result))
}

func TestHandler_ListForSelection_WithSubSubjects(t *testing.T) {
	handler, cleanup := setupSubjectsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/subjects/for-selection", nil)
	w := httptest.NewRecorder()

	handler.ListForSelection(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []SubjectForSelection
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)

	// Check if we have both main subjects and sub-subjects
	var mainSubjects, subSubjects []SubjectForSelection
	for _, subject := range result {
		if subject.ParentID == nil {
			mainSubjects = append(mainSubjects, subject)
		} else {
			subSubjects = append(subSubjects, subject)
		}
	}

	// Should have at least some main subjects
	assert.GreaterOrEqual(t, len(mainSubjects), 1)

	// If there are sub-subjects, they should have parent IDs
	for _, subSubject := range subSubjects {
		assert.NotNil(t, subSubject.ParentID)
		assert.NotEqual(t, uuid.Nil, *subSubject.ParentID)
	}
}

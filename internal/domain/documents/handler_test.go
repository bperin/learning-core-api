package documents

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/http/authz"
)

func setupHandler(t *testing.T) *Handler {
	t.Helper()
	// For testing, we use a nil GCS service since we're testing handler validation logic
	// In production, the GCS service would be properly initialized
	return NewHandler(nil, nil)
}

func TestHandler_GetSignedUploadURLAdmin(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"filename":     "document.pdf",
		"content_type": "application/pdf",
		"ttl_seconds":  600,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Add admin role to context
	ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result GetSignedUploadURLResponse
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.SignedURL)
	assert.Equal(t, "document.pdf", result.Filename)
	assert.NotEmpty(t, result.ExpiresAt)
}

func TestHandler_GetSignedUploadURLTeacher(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"filename":     "lesson.docx",
		"content_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Add teacher role to context
	ctx := authz.WithAuth(req.Context(), "teacher-456", []string{authz.RoleTeacher}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result GetSignedUploadURLResponse
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.SignedURL)
	assert.Equal(t, "lesson.docx", result.Filename)
}

func TestHandler_GetSignedUploadURLMissingFilename(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"content_type": "application/pdf",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetSignedUploadURLMissingContentType(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"filename": "document.pdf",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetSignedUploadURLInvalidContentType(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"filename":     "malicious.exe",
		"content_type": "application/x-msdownload",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetSignedUploadURLAllowedContentTypes(t *testing.T) {
	handler := setupHandler(t)

	allowedTypes := []string{
		"application/pdf",
		"text/plain",
		"text/markdown",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/msword",
		"text/html",
		"application/json",
	}

	for _, contentType := range allowedTypes {
		t.Run(contentType, func(t *testing.T) {
			payload := map[string]interface{}{
				"filename":     "document.ext",
				"content_type": contentType,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.GetSignedUploadURL(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "content type %s should be allowed", contentType)
		})
	}
}

func TestHandler_GetSignedUploadURLWithCustomTTL(t *testing.T) {
	handler := setupHandler(t)

	payload := map[string]interface{}{
		"filename":     "document.pdf",
		"content_type": "application/pdf",
		"ttl_seconds":  3600, // 1 hour
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/documents/signed-url", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := authz.WithAuth(req.Context(), "user-123", []string{authz.RoleAdmin}, []string{"read", "write"})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetSignedUploadURL(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result GetSignedUploadURLResponse
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.SignedURL)
	assert.NotEmpty(t, result.ExpiresAt)
}
